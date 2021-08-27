package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/internals/notifier"
	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/eventstream"

	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	groupID, ok := os.LookupEnv("NOTIFICATION_GROUP_ID")
	if !ok {
		groupID = "NOTIFICATION"
	}

	consumer, err := eventstream.GetKafkaConsumer(groupID, "demo")

	if err != nil {
		log.Println(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(1000)

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				n, err := notifier.GetNotifier(e)
				if err != nil {
					log.Printf("Error %v\n", err)
				}

				n.Send()
			case kafka.Error:
				log.Printf("ErrorL %v\n", e)
			case *kafka.Stats:
				var stats map[string]interface{}
				_ = json.Unmarshal([]byte(e.String()), &stats)
				log.Printf("Stats: %v messages (%v bytes) messages consumed\n", stats["rxmsgs"], stats["rxmsg_bytes"])
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}
}
