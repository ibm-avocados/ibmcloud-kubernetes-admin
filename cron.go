package main

import (
	"log"
	"time"

	"github.com/moficodes/ibmcloud-kubernetes-admin/ibmcloud"
)

var ticker *time.Ticker
var quit chan struct{}
var count int

func init() {
	ticker = time.NewTicker(5 * time.Second)
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

		schedules, err := session.GetDocument(accountID)
		if err != nil {
			log.Println(err)
		}

		for _, schedule := range schedules {
			log.Println(schedule.Status)
		}
	}
}
