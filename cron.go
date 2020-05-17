package main

import (
	"log"
	"strings"
	"time"

	"github.com/moficodes/ibmcloud-kubernetes-admin/ibmcloud"
)

var ticker *time.Ticker
var quit chan struct{}
var count int

func init() {
	ticker = time.NewTicker(20 * time.Second)
	quit = make(chan struct{})
	count = 0
}

func cron() {
	for {
		select {
		case <-ticker.C:
			// do stuff
			count++
			log.Printf("Timer called %d times", count)
			checkCloudant()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func runCron() {

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
		schedules, err := session.GetDocument(accountID)
		if err != nil {
			log.Println(err)
		}
		log.Println("schedule found : ", len(schedules))

		for _, schedule := range schedules {
			if schedule.Status == "scheduled" {
				// deal with creating the clusters and updating the schedule to created

				// get tags out of the schedule
				tags := strings.Split(schedule.Tags, ",")
				// for each cluster loop through and create cluster, ignore error for now.
				for _, createRequest := range schedule.ClusterRequests {
					response, err := session.CreateCluster(createRequest)
					if err != nil {
						log.Println("error creating cluster. investigate : ", createRequest.ClusterRequest.Name, err)
						continue
					}
					log.Println("created cluster :", response.ID)
					for _, tag := range tags {
						_, err := session.SetClusterTag(tag, response.ID, createRequest.ResourceGroup)
						if err != nil {
							log.Println("error setting tag : investigate ", createRequest.ClusterRequest.Name, err)
							continue
						}
						log.Println("created tag ", tag)
					}
				}
				schedule.Status = "created"
				log.Println("schedule id", schedule.ID, "schedule rev", schedule.Rev)
				if err := session.UpdateDocument(accountID, schedule.ID, schedule.Rev, schedule); err != nil {
					log.Println("could not update document", err)
					continue
				}
				log.Println("updated schedule")
			} else if schedule.Status == "created" {
				// deal with deleting the clusters and updating the schedule to completed
			} else {
				// idk what can be coming in this code block, since those are the only two status we check
			}
		}
	}
}
