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

func timed() {
	for {
		select {
		case <-ticker.C:
			// do stuff
			count++
			log.Printf("Timer called %d times", count)
			doStuff()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func doStuff() {
	err := ibmcloud.SetupAccount("test")
	log.Println(err)
	dbs, _ := ibmcloud.GetAllDbs()
	for _, db := range dbs {
		key, _ := ibmcloud.GetAPIKey(db)
		log.Printf("DB : %s, Key : %s", db, key)
	}
}
