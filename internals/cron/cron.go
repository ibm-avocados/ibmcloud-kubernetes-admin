package cron

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

var ticker *time.Ticker
var quit chan struct{}
var count int

func init() {
	_period := os.Getenv("TICKER_PERIOD")
	period, err := strconv.Atoi(_period)
	if err != nil || period == 0 {
		period = 3600
	}
	log.Printf("ticker running in %d seconds interval\n", period)
	ticker = time.NewTicker(time.Duration(period) * time.Second)
	quit = make(chan struct{})
	count = 0
}

func Start() {
	go cron()
}

func cron() {
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
		log.Println("error getting accounts")
	}

	for _, accountID := range accounts {
		session, err := ibmcloud.GetSessionFromCloudant(accountID)
		if err != nil {
			log.Println(err)
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
			if schedule.Status == "scheduled" {
				// deal with creating the clusters and updating the schedule to created
				log.Printf("creating %d clusters", count)
				// get tags out of the schedule
				tags := strings.Split(schedule.Tags, ",")

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
			} else if schedule.Status == "created" {
				// deal with deleting the clusters and updating the schedule to completed
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
