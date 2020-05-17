package ibmcloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/go-cloudant"
	"github.com/mitchellh/mapstructure"
)

var username string
var password string
var host string
var cclient *cloudant.Client

func SetupCloudant() {
	username = os.Getenv("CLOUDANT_USER_NAME")
	password = os.Getenv("CLOUDANT_PASSWORD")
	host = os.Getenv("CLOUDANT_HOST")
	var err error
	cclient, err = cloudant.NewClient(username, password)
	if err != nil {
		log.Println("cloudant password not working")
	}
	log.Println("cloudant setup")
}

func SetupAccount(accountID string) error {
	dbName := "db-" + accountID
	return setupDB(dbName)
}

func setupDB(dbName string) error {
	_, err := cclient.EnsureDB(dbName)
	return err
}

func checkExistingAPIKey(dbName string) error {
	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return err
	}

	return checkAPIKey(apiKey.APIKey, dbName)
}

func CheckAPIKey(accountID string) error {
	dbName := "db-" + accountID
	return checkExistingAPIKey(dbName)
}

func checkAPIKey(apiKey, dbName string) error {
	session, err := IAMAuthenticate(apiKey)
	if err != nil {
		return fmt.Errorf("not a valid api key %v", err)
	}

	// can we get the account name currently selected?

	accountID := strings.TrimPrefix(dbName, "db-")
	accounts, err := session.GetAccounts()
	found := false
	for _, account := range accounts.Resources {
		if account.Metadata.GUID == accountID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("api key not valid for current account")
	}
	return nil
}

func SetAPIKey(apiKey, accountID string) error {
	dbName := "db-" + accountID
	if err := checkAPIKey(apiKey, dbName); err != nil {
		return err
	}
	return setAPIKey(apiKey, dbName)
}

func setAPIKey(apiKey, dbName string) error {
	db := getDB(dbName)
	type apiDoc struct {
		ID     string `json:"_id"`
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

func UpdateAPIKey(apiKey, accountID string) error {
	dbName := "db-" + accountID
	if err := checkAPIKey(apiKey, dbName); err != nil {
		return err
	}
	return updateAPIKey(apiKey, dbName)
}

func updateAPIKey(newKey, dbName string) error {
	db := getDB(dbName)

	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return err
	}

	apiKey.APIKey = newKey

	newRev, err := db.UpdateDocument(apiKey.ID, apiKey.Rev, apiKey)
	if err != nil {
		return err
	}

	log.Printf("updated api key with new rev %s\n", newRev)
	return nil
}

func DeleteAPIKey(accountID string) error {
	dbName := "db-" + accountID
	return deleteAPIKey(dbName)
}

func deleteAPIKey(dbName string) error {
	db := getDB(dbName)

	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return err
	}

	newRev, err := db.DeleteDocument(apiKey.ID, apiKey.Rev)
	if err != nil {
		return err
	}
	log.Printf("document deleted rev %s\n", newRev)
	return nil
}

func GetDocument(accountID string) ([]Schedule, error) {
	dbName := "db-" + accountID
	return getDocument(dbName)
}

func getDocument(dbName string) ([]Schedule, error) {
	db := getDB(dbName)

	timeGT := time.Now().Unix()
	timeLT := time.Now().Add(time.Hour * 2).Unix()

	createQuery := cloudant.Query{}
	createQuery.Selector = make(map[string]interface{})
	createQuery.Selector["createAt"] = map[string]int64{
		"$gt": timeGT,
		"$lt": timeLT,
	}
	createQuery.Selector["status"] = map[string]string{
		"$eq": "scheduled",
	}

	resCreate, err := db.SearchDocument(createQuery)
	if err != nil {
		return nil, err
	}

	destroyQuery := cloudant.Query{}
	destroyQuery.Selector = make(map[string]interface{})
	destroyQuery.Selector["destroyAt"] = map[string]int64{
		"$gt": timeGT,
		"$lt": timeLT,
	}
	destroyQuery.Selector["status"] = map[string]string{
		"$eq": "created",
	}
	resDestroy, err := db.SearchDocument(destroyQuery)
	if err != nil {
		return nil, err
	}

	resJoin := append(resCreate, resDestroy...)

	res := make([]Schedule, len(resJoin))
	for i, elem := range resJoin {
		sched, ok := elem.(map[string]interface{})
		log.Println("sched id", sched["_id"])
		if !ok {
			log.Println("could not convert to type")
			return nil, err
		}
		var schedule Schedule
		if err := mapstructure.Decode(sched, &schedule); err != nil {
			log.Println("nothing is working")
		}
		res[i] = schedule
	}

	return res, nil
}

func CreateDocument(accountID string, data interface{}) error {
	dbName := "db-" + accountID
	return createDocument(dbName, data)
}

func createDocument(dbName string, data interface{}) error {
	db := getDB(dbName)

	id, rev, err := db.CreateDocument(data)
	if err != nil {
		return err
	}
	log.Printf("document set with id %s, rev %s\n", id, rev)
	return nil
}

func UpdateDocument(accountID, id, rev string, data interface{}) error {
	dbName := "db-" + accountID
	return updateDocument(dbName, id, rev, data)
}

func updateDocument(dbName, id, rev string, data interface{}) error {
	db := getDB(dbName)

	newRev, err := db.UpdateDocument(id, rev, data)
	if err != nil {
		log.Println("error here", err)
		return err
	}
	log.Printf("document updated with rev %s\n", newRev)
	return nil
}

func DeleteDocument(accountID, id, rev string) error {
	dbName := "db-" + accountID
	return deleteDocument(dbName, id, rev)
}

func deleteDocument(dbName, id, rev string) error {
	db := getDB(dbName)

	newRev, err := db.DeleteDocument(id, rev)
	if err != nil {
		return nil
	}
	log.Printf("document delete with rev %s\n", newRev)
	return nil
}

func GetSessionFromCloudant(accountID string) (*Session, error) {
	dbName := "db-" + accountID
	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return nil, err
	}

	session, err := IAMAuthenticate(apiKey.APIKey)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func updateDocumentDirect(dbName, id, rev string, doc interface{}) error {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
		"Content-Type":  "application/json",
	}

	json, err := json.Marshal(doc)
	if err != nil {
		log.Println("marhsal error")
		return err
	}
	var result map[string]interface{}
	url := fmt.Sprintf("https://%s/%s/%s", host, dbName, id)

	log.Println(url)

	err = put(url, header, nil, json, &result)
	if err != nil {
		log.Println("put error")
		return err
	}
	log.Println(result)
	return nil
}

func getAPIKey(dbName string) (*ApiKey, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
	}
	var result ApiKey
	url := fmt.Sprintf("https://%s/%s/api_key", host, dbName)

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetAllAccountIDs() ([]string, error) {
	result, err := getAllDbs()
	if err != nil {
		return nil, err
	}
	for i, val := range result {
		result[i] = strings.TrimPrefix(val, "db-")
	}
	return result, nil
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

func AddSchedule() error {
	return addSchedule("")
}

func addSchedule(dbName string) error {
	_ = getDB(dbName)
	return nil
}

func getDB(dbName string) *cloudant.DB {
	db, _ := cclient.EnsureDB(dbName)

	return db
}

/*
GetAPIKey


*/
