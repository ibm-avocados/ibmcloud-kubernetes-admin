package cron

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

type ScheduleError struct {
	Error   error
	Message string
}

type EmailData struct {
	Schedule ibmcloud.Schedule
	Errors   []ScheduleError
}

// Start
func Start() {
	_period := os.Getenv("TICKER_PERIOD")
	period, err := strconv.Atoi(_period)
	if err != nil || period == 0 {
		period = 3600
	}
	log.Printf("ticker running in %d seconds interval\n", period)
	ticker := time.NewTicker(time.Duration(period) * time.Second)
	quit := make(chan struct{})
	count := 0
	for {
		select {
		case <-ticker.C:
			// do stuff
			count++
			log.Printf("Timer called %d times", count)
			runCron()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func runCron() {
	checkCloudant()
}

func findMatchingVlan(private, public []ibmcloud.Vlan) (string, string) {
	rand.Seed(time.Now().UnixNano())
	pairs := make([][]string, 0)
	// for all the privateVlan find their matchin public vlan
	for _, privateVlan := range private {
		privateRouter := privateVlan.Properties.PrimaryRouter
		privateMatch := privateRouter[1:]
		for _, publicVlan := range public {
			publicRouter := publicVlan.Properties.PrimaryRouter
			publicMatch := publicRouter[1:]
			if privateMatch == publicMatch {
				match := []string{privateVlan.ID, publicVlan.ID}
				pairs = append(pairs, match)
			}
		}
	}
	// if their is no matching vlan available
	// return empty
	if len(pairs) == 0 {
		return "", ""
	}
	// return one of the pairs at random
	// this is to prevent overloading of a single vlan in a region
	// if there are a fiew vlan lets round robin this
	randomPair := pairs[rand.Intn(len(pairs))]
	return randomPair[0], randomPair[1]
}

func checkCloudant() {
	accounts, err := ibmcloud.GetAllAccountIDs()
	if err != nil {
		// basically means cloudant is not there or can not be connected to
		// no way to recover
		// only sane option is to contact admin
		if err := notification.EmailAdmin("Cloudant Not Available", "<strong>Check cloudant database</strong>"); err != nil {
			log.Println(err)
		}
		log.Println("error getting accounts")
	}

	for _, accountID := range accounts {
		session, err := ibmcloud.GetSessionFromCloudant(accountID)
		if err != nil {
			// could not get session
			// means either there was not api key! eek
			// or the api key was deleted need to notify account admins
			log.Println(err)
			notification.EmailAdmin("API key invalid/unavailable", fmt.Sprintf("<p>Check api key for %s</p>", accountID))
			continue
		}

		adminEmails, err := session.GetAccountAdminEmails(accountID)
		if err != nil || len(adminEmails) == 0 {
			notification.EmailAdmin("No account email available", "<h1>No account email available</h1>")
			continue
		}

		log.Println("checking schedules for account : ", accountID)
		schedules, err := session.GetDocument(accountID)
		if err != nil {
			log.Println(err)
		}
		log.Println("schedule found : ", len(schedules))

		for _, schedule := range schedules {
			notifyEmails := schedule.NotifyEmails
			if notifyEmails == nil || len(notifyEmails) == 0 {
				notifyEmails = adminEmails
			}

			emailData := EmailData{
				Errors: make([]ScheduleError, 0),
			}

			hasErrors := false

			count, err := strconv.Atoi(schedule.Count)
			if err != nil {
				log.Println("error converting count", schedule.Count, err)
				continue
			}

			name := schedule.ScheduleRequest.ScheduleRequest.Name
			if schedule.Status == "created" {
				log.Printf("deleting %d clusters", count)
				if count == len(schedule.Clusters) {
					for _, cluster := range schedule.Clusters {
						if err := session.DeleteCluster(cluster, schedule.ScheduleRequest.ResourceGroup, "true"); err != nil {
							log.Println("error deleting cluster, investigate : ", schedule.ScheduleRequest.ScheduleRequest.Name, err)
							hasErrors = true
							schedError := ScheduleError{
								Error:   err,
								Message: fmt.Sprintf("Error deleting cluster %s", cluster),
							}
							emailData.Errors = append(emailData.Errors, schedError)
							continue
						}
						log.Println("deleted cluster : ", cluster)
					}
				} else {
					for i := 1; i <= count; i++ {
						suffix := fmt.Sprintf("-%03d", i)
						clusterName := name + suffix
						if err := session.DeleteCluster(clusterName, schedule.ScheduleRequest.ResourceGroup, "true"); err != nil {
							log.Println("error deleting cluster, investigate : ", clusterName, err)
							hasErrors = true
							schedError := ScheduleError{
								Error:   err,
								Message: fmt.Sprintf("Error deleting cluster %s", clusterName),
							}
							emailData.Errors = append(emailData.Errors, schedError)
							continue
						}
						log.Println("deleted cluster : ", clusterName)
					}
				}
				schedule.Status = "completed"
				if hasErrors {
					schedule.Status = "delete-incomplete"
				}
			} else if schedule.Status == "scheduled" {
				// deal with creating the clusters and updating the schedule to created
				log.Printf("creating %d clusters", count)
				// get tags out of the schedule
				tags := strings.Split(schedule.Tags, ",")

				// for each cluster loop through and create cluster, ignore error for now.
				for i := 1; i <= count; i++ {
					datacenters := schedule.ScheduleRequest.ScheduleRequest.DataCenters
					datacenter := datacenters[i%len(datacenters)]
					vlans, err := session.GetDatacenterVlan(datacenter)
					if err != nil {
						// could not get vlan
						// skip the scheduling
						log.Println(err)
						continue
					}

					privateVlans := make([]ibmcloud.Vlan, 0)
					publicVlans := make([]ibmcloud.Vlan, 0)

					for _, vlan := range vlans {
						if vlan.Type == "private" {
							privateVlans = append(privateVlans, vlan)
						} else if vlan.Type == "public" {
							publicVlans = append(publicVlans, vlan)
						}
					}

					// at these point, ideally we have a list of private and public vlans
					// if theres nothing in this list
					// vlans got deleted
					// and we can no longer create the clusters (?) [might be ok to set empty]
					// at this situation email the account admins
					if len(privateVlans) == 0 || len(publicVlans) == 0 {
						notification.Email("No vlan available", "<h1>Add vlan for region</h1>", notifyEmails...)
						continue
					}

					privateVlan, publicVlan := findMatchingVlan(privateVlans, publicVlans)

					if privateVlan == "" || publicVlan == "" {
						notification.Email("No matching vlan available", "<h1>Add vlan with same router for region</h1>", notifyEmails...)
						continue
					}

					suffix := fmt.Sprintf("-%03d", i)

					var clusterRequest ibmcloud.ClusterRequest
					copier.Copy(&clusterRequest, &schedule.ScheduleRequest.ScheduleRequest)

					var createRequest ibmcloud.CreateClusterRequest
					createRequest.ClusterRequest = clusterRequest
					// TODO: set stuff
					createRequest.ClusterRequest.DataCenter = ""
					createRequest.ClusterRequest.PublicVlan = publicVlan
					createRequest.ClusterRequest.PrivateVlan = privateVlan
					createRequest.ClusterRequest.Name = name + suffix
					log.Println("trying to create cluster with name : ", createRequest.ClusterRequest.Name)
					response, err := session.CreateCluster(createRequest)
					if err != nil {
						log.Println("error creating cluster. investigate : ", createRequest.ClusterRequest.Name, err)
						hasErrors = true
						schedError := ScheduleError{
							Error:   err,
							Message: fmt.Sprintf("Error creting cluster %s", createRequest.ClusterRequest.Name),
						}
						emailData.Errors = append(emailData.Errors, schedError)
						continue
					}

					log.Println("created cluster. ID : ", response.ID, " Name : ", createRequest.ClusterRequest.Name)

					schedule.Clusters = append(schedule.Clusters, response.ID)

					for _, tag := range tags {
						_, err := session.SetClusterTag(tag, response.ID, schedule.ScheduleRequest.ResourceGroup)
						if err != nil {
							log.Println("error setting tag : investigate ", schedule.ScheduleRequest.ScheduleRequest.Name, err)
							hasErrors = true
							schedError := ScheduleError{
								Error:   err,
								Message: fmt.Sprintf("Error creting tag %s for cluster %s", tag, schedule.ScheduleRequest.ScheduleRequest.Name),
							}
							emailData.Errors = append(emailData.Errors, schedError)
							continue
						}
						log.Println("created tag ", tag)
					}
				}

				schedule.Status = "created"
				if hasErrors {
					schedule.Status = "create-incomplete"
				}
			} else {
				// idk what can be coming in this code block, since those are the only two status we check
			}
			if err := session.UpdateDocument(accountID, schedule.ID, schedule.Rev, schedule); err != nil {
				log.Println("could not update document", err)
				continue
			}
			log.Println("updated schedule status to ", schedule.Status)
			// at this point we should be able to email
			emailData.Schedule = schedule
			emailBody, err := getEmailBody(emailData)
			if err != nil {
				log.Println("could not get email body")
			}
			log.Println("will try send notification emails to : ", notifyEmails)
			if err := notification.Email("IBMCloud Kubernetes Admin Schedule executed", emailBody, notifyEmails...); err != nil {
				log.Println("error sending email")
			}

			// if its a workshop deploy cloud foundry and update github issue
			if !schedule.IsWorkshop {
				log.Println("Not a workshop")
				continue
			}

			apikey, err := session.GetAPIKey(accountID)
			if err != nil {
				log.Println("could not get api key")
				continue
			}
			metadata, err := session.GetAccountMetaData(accountID)
			if err != nil {
				log.Println("could not get account metadata")
				continue
			}
			setEnvs(accountID, apikey, metadata, schedule)
			// if the schedule is created we deploy the app
			// probably will deploy even if there was minor errors
			if schedule.Status == "created" {
				if err := deploy(apikey, metadata, schedule); err != nil {
					log.Println("could not deploy cloud foundry application")
					continue
				}

				if err := createComment(schedule, metadata, "templates/message.gotmpl"); err != nil {
					log.Println("could not update comment on github")
					continue
				}
			} else if schedule.Status == "completed" {
				if err := cleanUp(apikey, metadata, schedule); err != nil {
					log.Println("could not cleanup cloud foundry application")
					continue
				}
			}
		}
	}
}
