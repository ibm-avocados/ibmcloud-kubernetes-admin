package ibmcloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/IBM-Cloud/go-cloudant"
)

var username string
var password string
var host string
var cclient *cloudant.Client

func init() {
	username = os.Getenv("CLOUDANT_USER_NAME")
	password = os.Getenv("CLOUDANT_PASSWORD")
	host = username + ".cloudantnosqldb.appdomain.cloud"
	var err error
	cclient, err = cloudant.NewClient(username, password)
	if err != nil {
		log.Println("cloudant password not working")
	}
}

func SetupAccount(accountID string) error {
	return setupDB(accountID)
}

func setupDB(dbName string) error {
	_, err := cclient.CreateDB(dbName)
	return err
}

func SetAPIKey(apiKey, accountID string) error {
	return setAPIKey(apiKey, accountID)
}

func setAPIKey(apiKey, dbName string) error {
	db := cclient.DB(dbName)

	// check validity of api key
	// can we get a session?
	session, err := IAMAuthenticate(apiKey)
	if err != nil {
		return fmt.Errorf("not a valid api key ", err)
	}

	// can we get the account name currently selected?
	accounts, err := session.GetAccounts()
	found := false
	for _, account := range accounts.Resources {
		if account.Metadata.GUID == dbName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("api key not valid for current account")
	}

	type apiDoc struct {
		ID     string `json:"id"`
		APIKey string `json:"apikey"`
	}
	api := &apiDoc{
		ID:     "api_key",
		APIKey: apiKey,
	}
	id, rev, err := db.CreateDocument(api)
	if err != nil {
		return err
	}
	log.Println(id, rev, "api key set")
	return nil
}

func GetAPIKey(accountID string) (string, error) {
	return getAPIKey(accountID)
}

func getAPIKey(dbName string) (string, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
	}
	var result ApiKey
	url := fmt.Sprintf("https://%s/%s/api_key", host, dbName)

	err := fetch(url, header, nil, &result)
	if err != nil {
		return "", err
	}

	return result.APIKey, nil
}

func GetAllDbs() ([]string, error) {
	return getAllDbs()
}

func getAllDbs() ([]string, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
	}
	var result []string
	url := fmt.Sprintf("https://%s/_all_dbs", host)

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
