package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vault"
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
