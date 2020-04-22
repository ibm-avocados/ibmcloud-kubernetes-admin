package main

import (
	"log"
	"time"
)

var ticker *time.Ticker
var quit chan struct{}
var count int

func init() {
	ticker = time.NewTicker(600 * time.Second)
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

}
