package cron

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

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

func checkCloudant() {
	accounts, err := ibmcloud.GetAllAccountIDs()
	if err != nil {
		// basically means cloudant is not there or can not be connected to
		// no way to recover
		// only sane option is to contact admin
		notification.EmailAdmin("Cloudant Not Available", "<strong>Check cloudant database</strong>")
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
			if err := notification.EmailAdmin("No account email available", "<h1>No account email available</h1>"); err != nil {
				log.Println(err)
			}
		}

		log.Println("checking schedules for account : ", accountID)
		schedules, err := session.GetDocumentV2(accountID)
		if err != nil {
			log.Println(err)
		}
		log.Println("schedule found : ", len(schedules))

		for _, schedule := range schedules {

			count, err := strconv.Atoi(schedule.Count)
			if err != nil {
				log.Println("error converting count", schedule.Count, err)
				continue
			}

			name := schedule.CreateRequest.ClusterRequest.Name
			if schedule.Status == "created" {
				log.Printf("deleting %d clusters", count)
				if count == len(schedule.Clusters) {
					for _, cluster := range schedule.Clusters {
						if err := session.DeleteCluster(cluster, schedule.CreateRequest.ResourceGroup, "true"); err != nil {
							log.Println("error deleting cluster, investigate : ", schedule.CreateRequest.ClusterRequest.Name, err)
							continue
						}
						log.Println("deleted cluster : ", cluster)
					}
				} else {
					for i := 1; i <= count; i++ {
						suffix := fmt.Sprintf("-%03d", i)
						clusterName := name + suffix
						if err := session.DeleteCluster(clusterName, schedule.CreateRequest.ResourceGroup, "true"); err != nil {
							log.Println("error deleting cluster, investigate : ", clusterName, err)
							continue
						}
						log.Println("deleted cluster : ", clusterName)
					}
				}
				schedule.Status = "completed"
			} else if schedule.Status == "scheduled" {
				// deal with creating the clusters and updating the schedule to created
				log.Printf("creating %d clusters", count)
				// get tags out of the schedule
				tags := strings.Split(schedule.Tags, ",")

				vlans, err := session.GetDatacenterVlan(schedule.CreateRequest.ClusterRequest.DataCenter)
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
					notification.Email("No vlan available", "<h1>Add vlan for region</h1>")
					continue
				}

				// for each cluster loop through and create cluster, ignore error for now.
				for i := 1; i <= count; i++ {
					suffix := fmt.Sprintf("-%03d", i)
					schedule.CreateRequest.ClusterRequest.Name = name + suffix
					response, err := session.CreateCluster(schedule.CreateRequest)
					if err != nil {
						log.Println("error creating cluster. investigate : ", schedule.CreateRequest.ClusterRequest.Name, err)
						continue
					}

					log.Println("created cluster :", response.ID)

					schedule.Clusters = append(schedule.Clusters, response.ID)

					for _, tag := range tags {
						_, err := session.SetClusterTag(tag, response.ID, schedule.CreateRequest.ResourceGroup)
						if err != nil {
							log.Println("error setting tag : investigate ", schedule.CreateRequest.ClusterRequest.Name, err)
							continue
						}
						log.Println("created tag ", tag)
					}
				}

				schedule.Status = "created"
			} else {
				// idk what can be coming in this code block, since those are the only two status we check
			}
			if err := session.UpdateDocument(accountID, schedule.ID, schedule.Rev, schedule); err != nil {
				log.Println("could not update document", err)
				continue
			}
			log.Println("updated schedule status to ", schedule.Status)
		}
	}
}
