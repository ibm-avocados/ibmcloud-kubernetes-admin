package ibmcloud

import (
	"encoding/base64"
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
	log.Println("cloudant setup complete")
}

func SetupAccount(accountID string) error {
	dbName := "db-" + accountID
	return setupDB(dbName)
}

func setupDB(dbName string) error {
	_, err := cclient.EnsureDB(dbName)
	return err
}

func GetAPIKey(accountID string) (string, error) {
	dbName := "db-" + accountID
	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return "", err
	}
	return apiKey.APIKey, nil
}

func CheckAPIKey(accountID string) error {
	dbName := "db-" + accountID
	return checkExistingAPIKey(dbName)
}

func checkExistingAPIKey(dbName string) error {
	apiKey, err := getAPIKey(dbName)
	if err != nil {
		return err
	}

	return checkAPIKey(apiKey.APIKey, dbName)
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

func CreateAccountMetadata(accountID, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	dbName := "db-" + accountID
	return createAccountMetaData(dbName, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken)
}

func createMetaDataDocs(dbName, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	db := getDB(dbName)

	metaData := struct {
		ID               string `json:"id"`
		Org              string `json:"org"`
		Space            string `json:"space"`
		Region           string `json:"region"`
		IssueRepo        string `json:"issue_repo"`
		GrantClusterRepo string `json:"grant_cluster_repo"`
		GithubUser       string `json:"github_user"`
		GithubToken      string `json:"github_token"`
	}{
		ID:               "metadata",
		Org:              org,
		Space:            space,
		Region:           region,
		IssueRepo:        issueRepo,
		GrantClusterRepo: grantClusterRepo,
		GithubUser:       githubUser,
		GithubToken:      githubToken,
	}

	id, rev, err := db.CreateDocument(metaData)
	if err != nil {
		return err
	}
	log.Println(id, rev, "admin emails set")
	return nil
}

func UpdateAccountMetadata(accountID, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	dbName := "db-" + accountID
	return updateAccountMetaData(dbName, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken)
}

func updateAccountMetaData(dbName, org, space, region, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	db := getDB(dbName)

	metadata, err := getAccountMetaData(dbName)
	if err != nil {
		return err
	}

	metadata.Org = org
	metadata.Space = space
	metadata.Region = region
	metadata.IssueRepo = issueRepo
	metadata.GrantClusterRepo = grantClusterRepo
	metadata.GithubUser = githubUser
	metadata.GithubToken = githubToken

	newRev, err := db.UpdateDocument(metadata.ID, metadata.Rev, metadata)
	if err != nil {
		return err
	}
	log.Println("metadata updated with rev : ", newRev)
	return nil
}

func GetAccountMetaData(accountID string) (*AccountMetaData, error) {
	dbName := "db-" + accountID
	return getAccountMetaData(dbName)
}

func getAccountMetaData(dbName string) (*AccountMetaData, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
	}
	var result AccountMetaData
	url := fmt.Sprintf("https://%s/%s/metadata", host, dbName)

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateAdminEmails(accountId string, emails ...string) error {
	dbName := "db-" + accountId
	return createAdminEmail(dbName, emails...)
}

func createAdminEmail(dbName string, email ...string) error {
	db := getDB(dbName)

	admins := struct {
		ID     string   `json:"_id"`
		Emails []string `json:"emails"`
	}{
		ID:     "admins",
		Emails: email,
	}
	id, rev, err := db.CreateDocument(admins)
	if err != nil {
		return err
	}
	log.Println(id, rev, "admin emails set")
	return nil
}

func RemoveAdminEmails(accountID string, emails ...string) error {
	dbName := "db-" + accountID
	return removeAccountAdminEmails(dbName, emails...)
}

func removeAccountAdminEmails(dbName string, emails ...string) error {
	db := getDB(dbName)
	admins, err := getAccountAdminEmails(dbName)
	if err != nil {
		return err
	}
	adminEmails := admins.Emails
	for _, email := range emails {
		idx := find(adminEmails, email)
		if idx > 0 {
			admins.Emails = removeIndex(adminEmails, idx)
		}
	}
	newRev, err := db.UpdateDocument(admins.ID, admins.Rev, admins)
	if err != nil {
		return err
	}

	log.Println("updated email with ", newRev)
	return nil
}

func AddAdminEmails(accountID string, email ...string) error {
	dbName := "db-" + accountID
	return addAccountAdminEmails(dbName, email...)
}

func addAccountAdminEmails(dbName string, email ...string) error {
	db := getDB(dbName)
	admins, err := getAccountAdminEmails(dbName)
	if err != nil {
		log.Println("could not get valid admin emails")
		return err
	}

	admins.Emails = append(admins.Emails, email...)

	newRev, err := db.UpdateDocument(admins.ID, admins.Rev, admins)
	if err != nil {
		log.Println("email document update failed", err)
		return err
	}

	log.Printf("updated admin emails with new rev %s\n", newRev)
	return nil
}

func DeleteAdminEmails(accountID string) error {
	dbName := "db-" + accountID
	return deleteAccountAdminEmails(dbName)
}

func deleteAccountAdminEmails(dbName string) error {
	db := getDB(dbName)
	admins, err := getAccountAdminEmails(dbName)
	if err != nil {
		return err
	}
	newRev, err := db.UpdateDocument(admins.ID, admins.Rev, []string{})
	if err != nil {
		return err
	}

	log.Printf("removed admin emails %s", newRev)
	return nil
}

func GetAccountAdminEmails(accountID string) ([]string, error) {
	dbName := "db-" + accountID
	admins, err := getAccountAdminEmails(dbName)
	if err != nil {
		return nil, err
	}
	return admins.Emails, nil
}

func getAccountAdminEmails(dbName string) (*AccountAdminEmails, error) {
	authEncoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	header := map[string]string{
		"Authorization": "Basic " + authEncoded,
	}
	var result AccountAdminEmails
	url := fmt.Sprintf("https://%s/%s/admins", host, dbName)

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetAllDocument(accountID string) ([]interface{}, error) {
	dbName := "db-" + accountID
	return getAllDocument(dbName)
}

func getAllDocument(dbName string) ([]interface{}, error) {
	db := getDB(dbName)
	q := cloudant.Query{}
	q.Selector = make(map[string]interface{})
	q.Selector["status"] = map[string]interface{}{
		"$gt": "",
	}
	res, err := db.SearchDocument(q)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetDocument(accountID string) ([]Schedule, error) {
	dbName := "db-" + accountID
	return getUpcomingDocument(dbName)
}

func getUpcomingDocument(dbName string) ([]Schedule, error) {
	db := getDB(dbName)

	createQuery := cloudant.Query{}
	createQuery.Selector = make(map[string]interface{})
	createQuery.Selector["createAt"] = map[string]int64{
		"$gt": time.Now().Unix(),
		"$lt": time.Now().Add(time.Hour * 2).Unix(),
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
		"$gt": time.Now().Unix(),
		"$lt": time.Now().Add(time.Hour * 2).Unix(),
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

// func GetDocument(accountID string) ([]Schedule, error) {
// 	dbName := "db-" + accountID
// 	return getUpcomingDocument(dbName)
// }

// func getUpcomingDocument(dbName string) ([]Schedule, error) {
// 	db := getDB(dbName)

// 	createQuery := cloudant.Query{}
// 	createQuery.Selector = make(map[string]interface{})
// 	createQuery.Selector["createAt"] = map[string]int64{
// 		"$gt": time.Now().Unix(),
// 		"$lt": time.Now().Add(time.Hour * 2).Unix(),
// 	}
// 	createQuery.Selector["status"] = map[string]string{
// 		"$eq": "scheduled",
// 	}

// 	resCreate, err := db.SearchDocument(createQuery)
// 	if err != nil {
// 		return nil, err
// 	}

// 	destroyQuery := cloudant.Query{}
// 	destroyQuery.Selector = make(map[string]interface{})
// 	destroyQuery.Selector["destroyAt"] = map[string]int64{
// 		"$gt": time.Now().Unix(),
// 		"$lt": time.Now().Add(time.Hour * 2).Unix(),
// 	}
// 	destroyQuery.Selector["status"] = map[string]string{
// 		"$eq": "created",
// 	}
// 	resDestroy, err := db.SearchDocument(destroyQuery)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resJoin := append(resCreate, resDestroy...)

// 	res := make([]Schedule, len(resJoin))
// 	for i, elem := range resJoin {
// 		sched, ok := elem.(map[string]interface{})
// 		if !ok {
// 			log.Println("could not convert to type")
// 			return nil, err
// 		}
// 		var schedule Schedule
// 		if err := mapstructure.Decode(sched, &schedule); err != nil {
// 			log.Println("nothing is working")
// 		}
// 		res[i] = schedule
// 	}

// 	return res, nil
// }

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

	_, err := db.UpdateDocument(id, rev, data)
	if err != nil {
		log.Println("error here", err)
		return err
	}
	return nil
}

func DeleteDocument(accountID, id, rev string) error {
	dbName := "db-" + accountID
	return deleteDocument(dbName, id, rev)
}

func deleteDocument(dbName, id, rev string) error {
	db := getDB(dbName)

	_, err := db.DeleteDocument(id, rev)
	if err != nil {
		return nil
	}
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

func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

/*
GetAPIKey


*/
