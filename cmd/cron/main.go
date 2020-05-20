package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/cron"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func main() {
	ibmcloud.SetupCloudant()
	cron.Start()
}
