package main

import (
	"log"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra/ibmcloud"
	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/router/kubeadmin"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	provider := ibmcloud.NewProvider()
	router, err := kubeadmin.NewRouter(provider)
	if err != nil {
		return err
	}
	return router.Serve(":9030")
}
