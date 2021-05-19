package main

import (
	"log"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/infra/ibmcloud"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/router/kubeadmin"
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
