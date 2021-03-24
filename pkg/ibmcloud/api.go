package ibmcloud

// TODO: return errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// protocol
const protocol = "https://"

// subdomains
const (
	subdomainIAM                = "iam."
	subdomainUserManagement     = "user-management."
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
	identityEndpoint       = protocol + subdomainIAM + api + "/identity/.well-known/openid-configuration"
	userPreferenceEndpoint = protocol + "user-preferences.ng.bluemix.net/v1/users/"
	accountsEndpoint       = protocol + subdomainAccounts + api + "/coe/v2/accounts"
	resourcesEndpoint      = protocol + subdomainResourceController + api + "/v2/resource_instances"
	resourceKeysEndpoint   = protocol + subdomainResourceController + api + "/v2/resource_keys"
	containersEndpoint     = protocol + subdomainClusters + api + "/global/v1"
	usersEndpoint          = protocol + subdomainUsers + api + "/v2"
	tagEndpoint            = protocol + subdomainTags + api + "/v3/tags"
	billingEndpoint        = protocol + subdomainBilling + api + "/v4/accounts"
	resourceEndoint        = protocol + subdomainResourceController + api + "/v1/resource_groups"
	apikeyEndpoint         = protocol + subdomainIAM + api + "/v1/apikeys"
	iamEndpoint            = protocol + subdomainIAM + api + "/v2/groups"
	userManagementEndpoint = protocol + subdomainUserManagement + api + "/v2/accounts"
	policyEndpoint         = protocol + subdomainIAM + api + "/v1/policies"
)

const (
	clusterEndpoint     = containersEndpoint + "/clusters"
	versionEndpount     = containersEndpoint + "/versions"
	locationEndpoint    = containersEndpoint + "/locations"
	zonesEndpoint       = containersEndpoint + "/zones"
	datacentersEndpoint = containersEndpoint + "/datacenters"
)

// grant types
const (
	passcodeGrantType     = "urn:ibm:params:oauth:grant-type:passcode"
	apikeyGrantType       = "urn:ibm:params:oauth:grant-type:apikey"
	refreshTokenGrantType = "refresh_token"
)

var basicAuth = "Basic " + os.Getenv("IBM_LOGIN_USER")

//// useful for loagging
// bodyBytes, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	panic(err)
// }
// bodyString := string(bodyBytes)
// log.Println(bodyString)
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

func checkToken(token, apikey string) (*ApiKeyDetails, error) {
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"IAM-Apikey":    apikey,
	}

	endpoint := apikeyEndpoint + "/details"

	var res ApiKeyDetails
	err := fetch(endpoint, header, nil, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
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

func getUserInfo(endpoint string, token string) (*UserInfo, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint cannot be empty")
	}
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	var result UserInfo
	err := fetch(endpoint, header, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func getUserPreference(userID, token string) (*User, error) {
	endpoint := userPreferenceEndpoint + userID

	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	var result User
	err := fetch(endpoint, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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

func getZones(showFlavors, location string) ([]Zone, error) {
	var result []Zone
	query := map[string]string{
		"showFlavors": showFlavors,
	}
	if len(location) > 0 {
		query["location"] = location
	}
	err := fetch(containersEndpoint+"/zones", nil, query, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAccessGroups(token, accountID string) (*AccessGroups, error) {
	var result AccessGroups
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{
		"account_id": accountID,
	}

	err := fetch(iamEndpoint, header, query, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Add user to account
func inviteUserToAccount(token, accountID, email string) (*UserInviteList, error) {
	var result UserInviteList
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	usersToInvite := []UserInvite{UserInvite{Email: email, AccountRole: "Member"}}
	userInviteList := UserInviteList{Users: usersToInvite}

	body, err := json.Marshal(userInviteList)
	if err != nil {
		return nil, err
	}

	inviteUserEndpoint := userManagementEndpoint + "/" + accountID + "/users/"

	err = postBody(inviteUserEndpoint, header, nil, body, &result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

// Add user to access group
func addMemberToAccessGroup(token, accessGroupID, iamID, memberType string) (*MemberList, error) {
	var result MemberList
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accepts":       "application/json",
	}

	membersToAdd := []Member{Member{IamID: iamID, Type: memberType}}
	memberAddList := MemberList{membersToAdd}

	body, err := json.Marshal(memberAddList)
	if err != nil {
		return nil, err
	}

	addMemberEndpoint := iamEndpoint + "/" + accessGroupID + "/members"

	err = put(addMemberEndpoint, header, nil, body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreatePolicy
func createPolicy(token, accountID, iamID, serviceName, serviceInstance, role string) (*PolicyResponse, error) {
	var result PolicyResponse
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
		"Accepts":       "application/json",
	}

	policy := Policy{
		Type:        "access",
		Description: "Access to instance",
		Subjects: []AttributeList{
			{
				[]Attribute{
					Attribute{
						Name:  "iam_id",
						Value: iamID,
					},
				},
			},
		},
		Roles: []Role{
			Role{role},
		},
		Resources: []AttributeList{
			{
				[]Attribute{
					Attribute{
						Name:  "accountId",
						Value: accountID,
					},
					Attribute{
						Name:  "serviceName",
						Value: serviceName,
					},
					Attribute{
						Name:  "serviceInstance",
						Value: serviceInstance,
					},
				},
			},
		},
	}

	body, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	err = postBody(policyEndpoint, header, nil, body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func isMemberOfAccessGroup(token, accessGroupID, iamID string) error {

	header := map[string]string{
		"Authorization":         "Bearer " + token,
	}

	checkMembershipEndpoint := iamEndpoint + "/" + accessGroupID + "/members/" + iamID
	err := head(checkMembershipEndpoint, header, nil, nil)
	if err != nil {
		return err
	}
	log.Println("User: " + iamID + " is a member of " + accessGroupID)
	return nil
}

func getAccountResources(token, accountID string) (*AccountResources, error) {
	var result AccountResources
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}

	query := map[string]string{
		"account_id": accountID,
	}

	err := fetch(resourceEndoint, header, query, &result)
	if err != nil {
		return nil, err
	}
	//"/v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5"
	return &result, nil
}

func getDatacenterVlan(token, refreshToken, datacenter string) ([]Vlan, error) {
	var result []Vlan
	header := map[string]string{
		"Authorization":        "Bearer " + token,
		"X-Auth-Refresh-Token": refreshToken,
	}

	url := datacentersEndpoint + "/" + datacenter + "/vlans"

	err := fetch(url, header, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getVersions() (*ClusterVersion, error) {
	var result ClusterVersion
	err := fetch(versionEndpount, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func getLocations() ([]Location, error) {
	var result []Location
	err := fetch(locationEndpoint, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getGeoLocations(geo string) ([]Location, error) {
	locations, err := getLocations()
	if err != nil {
		return nil, err
	}

	geoLocations := make([]Location, 0, 10)

	for _, location := range locations {
		if location.Kind == "dc" && location.Geography == geo {
			geoLocations = append(geoLocations, location)
		}
	}
	return geoLocations, nil
}

func getMachineTypes(datacenter, serverType, os string, cpuLimit, memoryLimit int) ([]MachineFlavor, error) {
	var result []MachineFlavor
	machineTypeEndpoint := fmt.Sprintf("%s/%s/machine-types", datacentersEndpoint, datacenter)
	err := fetch(machineTypeEndpoint, nil, nil, &result)
	if err != nil {
		return nil, err
	}
	if serverType != "" && os != "" {
		filtered := make([]MachineFlavor, 0)
		toLower := strings.ToLower
		atoi := strconv.Atoi
		for _, machine := range result {
			cpu, _ := atoi(machine.Cores)
			memory, _ := atoi(strings.ReplaceAll(machine.Memory, "GB", ""))
			if toLower(machine.ServerType) == toLower(serverType) &&
				toLower(machine.Os) == toLower(os) &&
				cpu <= cpuLimit &&
				memory <= memoryLimit {
				filtered = append(filtered, machine)
			}
		}
		return filtered, nil
	}
	return result, nil
}

func getCluster(token, clusterID, resourceGroup string) (*Cluster, error) {
	var result Cluster
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": resourceGroup,
	}
	err := fetch(clusterEndpoint+"/"+clusterID, header, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func getClusters(token, location string) ([]*Cluster, error) {
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

	// wg := &sync.WaitGroup{}

	// for _, cluster := range result {
	// 	time.Sleep(10 * time.Millisecond)
	// 	wg.Add(1)
	// 	go func(cluster *Cluster) {
	// 		tags, err := getTags(token, cluster.Crn)
	// 		if err != nil {
	// 			log.Println("error for tag: ", cluster.Name)
	// 			log.Println("error : ", err)
	// 		} else {
	// 			cluster.Tags = make([]string, len(tags.Items))
	// 			for i, val := range tags.Items {
	// 				cluster.Tags[i] = val.Name
	// 			}
	// 		}
	// 		wg.Done()
	// 	}(cluster)
	// 	wg.Add(1)
	// 	go func(cluster *Cluster) {
	// 		workers, err := getClusterWorkers(token, cluster.ID)
	// 		if err != nil {
	// 			log.Println("error for worker: ", cluster.Name)
	// 			log.Println("error : ", err)
	// 		} else {
	// 			cluster.Workers = workers
	// 			cost, err := getBillingData(token, accountID, cluster.Crn, workers)
	// 			if err != nil {
	// 				log.Println("error for cost: ", cluster.Name)
	// 			}
	// 			cluster.Cost = cost
	// 		}
	// 		wg.Done()
	// 	}(cluster)
	// }

	// wg.Wait()
	return result, nil
}

func getBillingData(token, accountID, clusterID, resourceInstanceID string) (string, error) {
	currentMonth := time.Now().Format("2006-01")
	workers, err := getClusterWorkers(token, clusterID)
	if err != nil {
		return "N/A", err
	}
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

func createCluster(token string, request CreateClusterRequest) (*CreateClusterResponse, error) {
	var result CreateClusterResponse
	header := map[string]string{
		"Authorization":         "Bearer " + token,
		"X-Auth-Resource-Group": request.ResourceGroup,
	}

	body, err := json.Marshal(request.ClusterRequest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = postBody(clusterEndpoint, header, nil, body, &result)

	if err != nil {
		log.Println("error creating cluster : ", request.ClusterRequest.Name, err)
		return nil, err
	}
	log.Printf("cluster created. id :%s => name: %s", result.ID, request.ClusterRequest.Name)
	return &result, nil
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
	log.Println("cluster deleted, id :", id)
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

func setClusterTags(token, tag, clusterID, resourceGroup string) (*TagResult, error) {
	cluster, err := getCluster(token, clusterID, resourceGroup)
	if err != nil {
		log.Println("get cluster : ", err)
		return nil, err
	}
	crn := cluster.Crn

	resources := make([]Resource, 1)
	resources[0] = Resource{ResourceID: crn}
	updateTag := UpdateTag{TagName: tag, Resources: resources}
	tagResult, err := setTags(token, updateTag)
	if err != nil {
		log.Println("set tag : ", err)
		return nil, err
	}
	return tagResult, nil
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
