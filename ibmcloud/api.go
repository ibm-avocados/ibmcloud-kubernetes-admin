package ibmcloud

// TODO: return errors

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/certifi/gocertifi"
)

// protocol
const protocol = "https://"

// subdomains
const (
	subdomainIAM                = "iam."
	subdomainAccounts           = "accounts."
	subdomainResourceController = "resource-controller."
	subdomainClusters           = "containers."
	subdomainUsers              = "users."
	subdomainTags               = "tags.global-search-tagging."
	subdomainBilling            = "billing."
)

// domain
const api = "cloud.ibm.com"

// endpoints
const (
	identityEndpoint     = protocol + subdomainIAM + api + "/identity/.well-known/openid-configuration"
	accountsEndpoint     = protocol + subdomainAccounts + api + "/coe/v2/accounts"
	resourcesEndpoint    = protocol + subdomainResourceController + api + "/v2/resource_instances"
	resourceKeysEndpoint = protocol + subdomainResourceController + api + "/v2/resource_keys"
	containersEndpoint   = protocol + subdomainClusters + api + "/global/v1"
	usersEndpoint        = protocol + subdomainUsers + api + "/v2"
	clusterEndpoint      = protocol + subdomainClusters + api + "/global/v1/clusters"
	tagEndpoint          = protocol + subdomainTags + api + "/v3/tags"
	billingEndpoint      = protocol + subdomainBilling + api + "/v4/accounts"
)

// grant types
const (
	passcodeGrantType     = "urn:ibm:params:oauth:grant-type:passcode"
	apikeyGrantType       = "urn:ibm:params:oauth:grant-type:apikey"
	refreshTokenGrantType = "refresh_token"
)

const basicAuth = "Basic Yng6Yng="

// TODO: logical timeout, 10 seconds wasn't long enough.
var client *http.Client

func init() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	certPool, _ := gocertifi.CACerts()
	transport.TLSClientConfig = &tls.Config{RootCAs: certPool}
	client = &http.Client{Transport: transport, Timeout: time.Duration(30 * time.Second)}
}

//// useful for loagging
// bodyBytes, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	panic(err)
// }
// bodyString := string(bodyBytes)
// fmt.Println(bodyString)
////

func timeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	log.Printf("TIME: %s took %s\n", name, elapsed)
}

func getError(resp *http.Response) error {
	var errorTemplate ErrorMessage
	if err := json.NewDecoder(resp.Body).Decode(&errorTemplate); err != nil {
		return err
	}
	if errorTemplate.Error != nil {
		return errors.New(errorTemplate.Error[0].Message)
	}
	if errorTemplate.Errors != nil {
		return errors.New(errorTemplate.Errors[0].Message)
	}
	return errors.New("unknown")
}

func getIdentityEndpoints() (*IdentityEndpoints, error) {
	result := &IdentityEndpoints{}
	err := fetch(identityEndpoint, nil, nil, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getToken(endpoint string, otp string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", passcodeGrantType)
	form.Add("passcode", otp)

	result := Token{}
	err := postForm(endpoint, header, nil, form, &result)

	if err != nil {
		log.Println("error in post form")
		return nil, err
	}

	return &result, nil
}

func getTokenFromIAM(endpoint string, apikey string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", apikeyGrantType)
	form.Add("apikey", apikey)

	result := &Token{}
	err := postForm(endpoint, header, nil, form, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func upgradeToken(endpoint string, refreshToken string, accountID string) (*Token, error) {
	header := map[string]string{
		"Authorization": basicAuth,
	}

	form := url.Values{}
	form.Add("grant_type", refreshTokenGrantType)
	form.Add("refresh_token", refreshToken)
	if accountID != "" {
		form.Add("bss_account", accountID)
	}

	result := &Token{}
	err := postForm(endpoint, header, nil, form, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getAccounts(endpoint *string, token string) (*Accounts, error) {
	if endpoint == nil {
		endpointString := accountsEndpoint
		endpoint = &endpointString
	} else {
		endpointString := accountsEndpoint + *endpoint
		endpoint = &endpointString
	}

	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	var result Accounts
	err := fetch(*endpoint, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getZones() ([]Zone, error) {
	var result []Zone
	header := map[string]string{}
	err := fetch(containersEndpoint+"/zones", header, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getLocations() ([]Location, error) {
	var result []Location
	err := fetch(containersEndpoint+"/zones", nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getClusters(token, accountID, location string) ([]*Cluster, error) {
	defer timeTaken(time.Now(), "GetCluster :")
	var result []*Cluster
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{}
	if len(location) > 0 {
		query["location"] = location
	}

	err := fetch(clusterEndpoint, header, query, &result)

	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}

	for _, cluster := range result {
		time.Sleep(10 * time.Millisecond)
		wg.Add(1)
		go func(cluster *Cluster) {
			tags, err := getTags(token, cluster.Crn)
			if err != nil {
				fmt.Println("error for tag: ", cluster.Name)
				fmt.Println("error : ", err)
			} else {
				cluster.Tags = make([]string, len(tags.Items))
				for i, val := range tags.Items {
					cluster.Tags[i] = val.Name
				}
			}
			wg.Done()
		}(cluster)
		wg.Add(1)
		go func(cluster *Cluster) {
			workers, err := getClusterWorkers(token, cluster.ID)
			if err != nil {
				fmt.Println("error for worker: ", cluster.Name)
				fmt.Println("error : ", err)
			} else {
				cluster.Workers = workers
				cost, err := getBillingData(token, accountID, cluster.Crn, workers)
				if err != nil {
					fmt.Println("error for cost: ", cluster.Name)
				}
				cluster.Cost = cost
			}
			wg.Done()
		}(cluster)
	}

	wg.Wait()
	return result, nil
}

func getBillingData(token, accountID, resourceInstanceID string, workers []Worker) (string, error) {
	currentMonth := time.Now().Format("2006-01")
	total := 0.0
	for _, worker := range workers {
		usage, err := getResourceUsagePerNode(token, accountID, currentMonth, resourceInstanceID, worker.ID)
		if err != nil {
			log.Printf("error getting resource usage %v\n", err)
			return "N/A", err
		}
		costForWorker := calcuateCostFromResourceUsage(usage)
		total += costForWorker
	}

	s := fmt.Sprintf("%.2f", total)
	return s, nil
}

func calcuateCostFromResourceUsage(usage *ResourceUsage) float64 {
	total := 0.0
	for _, resource := range usage.Resources {
		for _, use := range resource.Usage {
			total += use.Cost
		}
	}
	return total
}

func deleteCluster(token, id, resourceGroup, deleteResources string) error {
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": resourceGroup,
	}

	query := map[string]string{
		"deleteResources": deleteResources,
	}

	deleteEndpoint := clusterEndpoint + "/" + id
	err := delete(deleteEndpoint, header, query, nil)
	if err != nil {
		return err
	}
	return nil
}

func getClusterWorkers(token, id string) ([]Worker, error) {
	var result []Worker
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	workerEndpoint := clusterEndpoint + "/" + id + "/workers"

	err := fetch(workerEndpoint, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getResourceUsagePerNode(token, accountID, billingMonth, resourceInstanceID, workerID string) (*ResourceUsage, error) {
	var result ResourceUsage
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	crn := strings.ReplaceAll(resourceInstanceID, "::", ":worker:") + workerID
	query := map[string]string{
		"resource_id":          "containers-kubernetes",
		"_names":               "true",
		"resource_instance_id": crn,
	}

	endpoint := billingEndpoint + "/" + accountID + "/resource_instances/usage/" + billingMonth

	err := fetch(endpoint, header, query, &result)

	if err != nil {
		return nil, fmt.Errorf("error fetching resources usage %v", err)
	}

	return &result, err
}

func getTags(token string, crn string) (*Tags, error) {

	var result Tags
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	query := map[string]string{
		"attached_to": crn,
	}
	err := fetch(tagEndpoint, header, query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func setTags(token string, updateTag UpdateTag) (*TagResult, error) {
	setTagsEndpoint := tagEndpoint + "/" + "attach"
	return updateTags(setTagsEndpoint, token, updateTag)
}

func deleteTags(token string, updateTag UpdateTag) (*TagResult, error) {
	setTagsEndpoint := tagEndpoint + "/" + "detach"

	return updateTags(setTagsEndpoint, token, updateTag)
}

func updateTags(endpoint, token string, updateTag UpdateTag) (*TagResult, error) {
	var result TagResult
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	query := map[string]string{
		"providers": "ghost",
	}

	body, err := json.Marshal(updateTag)
	if err != nil {
		return nil, err
	}
	err = postBody(endpoint, header, query, body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
