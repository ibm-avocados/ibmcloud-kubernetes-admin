package main

import (
	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/vault"
	_ "github.com/joho/godotenv/autoload"
	"fmt"
	"log"
)

func main() {
	client, err := vault.GetVaultClient()
	if err != nil {
		log.Fatalln(err)
	}

	secret, err := client.Logical().Read("generic/user/mofizur-rahman/<accountID>")
	fmt.Printf("%+v", secret)
}
