package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/moficodes/ibmcloud-kubernetes-admin/ibmcloud"
)

const (
	errorMessageFormat = `{"msg": "error: %s"}`
	statusOkMessage    = `{"status": "ok"}`
	sessionName        = "cloud_session"
	accessToken        = "access_token"
	refreshToken       = "refresh_token"
	expiration         = "expiration"
	cookiePath         = "/api/v1"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	handleError(w, http.StatusNotFound, "not found")
}

func tokenEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	endpoints, err := ibmcloud.GetIdentityEndpoints()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get endpoints")
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(endpoints)
}

func resourceGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	accountResources, err := session.GetAccountResources(accountID)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not account resources", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(accountResources)
}

func vlanEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	vars := mux.Vars(r)

	datacenters, ok := vars["datacenter"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	vlans, err := session.GetDatacenterVlan(datacenters)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get vlans")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(vlans)
}

func versionEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	versions, err := ibmcloud.GetVersions()
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(versions)
}

func locationEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	locations, err := ibmcloud.GetLocations()
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(locations)
}

func locationGeoEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)

	geo, ok := vars["geo"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	locations, err := ibmcloud.GetGeoLocations(geo)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get locations")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(locations)
}

func zonesEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	showFlavors := r.FormValue("showFlavors")
	location := r.FormValue("location")
	zones, err := ibmcloud.GetZones(showFlavors, location)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not load zones")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(zones)
}

func machineTypeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	serverType := r.FormValue("type")
	os := r.FormValue("os")
	cpuLimitStr := r.FormValue("cpuLimit")
	cpuLimit, err := strconv.Atoi(cpuLimitStr)
	if err != nil {
		handleError(w, http.StatusBadRequest, "cpuLimit should be a number")
		return
	}
	memoryLimitStr := r.FormValue("memoryLimit")
	memoryLimit, err := strconv.Atoi(memoryLimitStr)
	if err != nil {
		handleError(w, http.StatusBadRequest, "memoryLimit should be a number")
		return
	}

	vars := mux.Vars(r)

	datacenter, ok := vars["datacenter"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}
	flavors, err := ibmcloud.GetMachineType(datacenter, serverType, os, cpuLimit, memoryLimit)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not load flavor")
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(flavors)
}

func authenticationWithAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	if !session.IsValid() {
		handleError(w, http.StatusUnauthorized, "session not valid")
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", body["id"])

	log.Println("Account id", accountID)

	accountSession, err := session.BindAccountToToken(accountID)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not bind account to token", err.Error())
		return
	}

	setCookie(w, accountSession)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not decode", err.Error())
		return
	}

	otp := fmt.Sprintf("%v", body["otp"])

	log.Println(otp)
	session, err := ibmcloud.Authenticate(otp)
	if err != nil {
		log.Println("could not authenticate with the otp provided")
		log.Println(err.Error())
		handleError(w, http.StatusUnauthorized, "could not authenticate with the otp provided", err.Error())
		return
	}

	log.Println(session.Token.Expiration)

	setCookie(w, session)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func setCookie(w http.ResponseWriter, session *ibmcloud.Session) {
	accessTokenCookie := http.Cookie{Name: accessToken, Value: session.Token.AccessToken, Path: cookiePath}
	http.SetCookie(w, &accessTokenCookie)

	refreshTokenCookie := http.Cookie{Name: refreshToken, Value: session.Token.RefreshToken, Path: cookiePath}
	http.SetCookie(w, &refreshTokenCookie)

	expirationStr := strconv.Itoa(session.Token.Expiration)

	expirationCookie := http.Cookie{Name: expiration, Value: expirationStr, Path: cookiePath}
	http.SetCookie(w, &expirationCookie)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	log.Println(session.Token.Expiration)

	if !session.IsValid() {
		handleError(w, http.StatusUnauthorized, "session expired")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func accountListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		log.Printf("could not get session %v\n", err)
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	accounts, err := session.GetAccounts()
	if err != nil {
		log.Printf("could not get accounts using access token %v\n", err)
		handleError(w, http.StatusUnauthorized, "could not get accounts using access token", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(accounts)
}

func clusterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	id := fmt.Sprintf("%v", body["id"])
	resoueceGroup := fmt.Sprintf("%v", body["resourceGroup"])
	deleteResources := fmt.Sprintf("%v", body["deleteResources"])
	log.Println(id, resoueceGroup, deleteResources)
	err = session.DeleteCluster(id, resoueceGroup, deleteResources)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func clusterCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.CreateClusterRequest

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode ", err.Error())
		return
	}

	createResponse, err := session.CreateCluster(body)

	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not create cluster", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(createResponse)
}

func clusterHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	resourceGroup, ok := body["resourceGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no tag attached to body", err.Error())
		return
	}
	clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)
	cluster, err := session.GetCluster(clusterID, clusterResourceGroup)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get cluster", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(cluster)
}

func clusterListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	clusters, err := session.GetClusters("")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get clusters", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(clusters)
}

func getBillingHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	crn, ok := body["crn"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no crn attached to body", err.Error())
		return
	}

	clusterCRN := fmt.Sprintf("%v", crn)

	accnt, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no account id attached to body", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", accnt)

	clustr, ok := body["clusterID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no cluster id attached to body", err.Error())
		return
	}

	clusterID := fmt.Sprintf("%v", clustr)

	billing, err := session.GetBillingData(accountID, clusterID, clusterCRN)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get billing info", err.Error())
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `{"bill": "%s"}`, billing)
}

func clusterWorkerListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}

	workers, err := session.GetWorkers(clusterID)

	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get workers for cluster : ", clusterID, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(workers)
}

func deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	res, err := session.DeleteTag(body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete tag", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, deleted %d tags"}`, len(res.Results))
}

func getTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	crn, ok := body["crn"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no crn attached to body", err.Error())
		return
	}

	clusterCRN := fmt.Sprintf("%v", crn)

	tags, err := session.GetTags(clusterCRN)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get tags", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(tags)
}

func setClusterTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	tag, ok := body["tag"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no tag attached to body", err.Error())
		return
	}

	resourceGroup, ok := body["resourceGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no resourceID attached to body", err.Error())
		return
	}

	vars := mux.Vars(r)

	clusterID, ok := vars["clusterID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get clusterID")
		return
	}
	clusterTag := fmt.Sprintf("%v", tag)
	clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)

	res, err := session.SetClusterTag(clusterTag, clusterID, clusterResourceGroup)
	//
	log.Println(err)
	if err != nil {
		log.Println("set cluster tag : ", err)
		handleError(w, http.StatusInternalServerError, "could not update tag", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, wrote %d tags"}`, len(res.Results))

}

func setTagHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	res, err := session.SetTag(body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not update tag", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "ok, wrote %d tags"}`, len(res.Results))
}

func setAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_accountID, ok := body["accountID"]
	if !ok {
		if err != nil {
			handleError(w, http.StatusBadRequest, "could not get account id", err.Error())
			return
		}
	}
	accountID := fmt.Sprintf("%v", _accountID)

	_apiKey, ok := body["apiKey"]
	if !ok {
		if err != nil {
			handleError(w, http.StatusBadRequest, "could not get apikey", err.Error())
			return
		}
	}
	apiKey := fmt.Sprintf("%v", _apiKey)

	if err := session.SetAPIKey(apiKey, accountID); err != nil {
		handleError(w, http.StatusUnauthorized, "could not set apikey", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

func checkAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_accountID, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "could not get account id")
		return
	}
	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.CheckAPIKey(accountID); err != nil {
		handleError(w, http.StatusUnauthorized, "error validating api key", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

func updateAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_accountID, ok := body["accountID"]
	if !ok {
		if err != nil {
			handleError(w, http.StatusBadRequest, "could not get account id", err.Error())
			return
		}
	}
	accountID := fmt.Sprintf("%v", _accountID)

	_apiKey, ok := body["apiKey"]
	if !ok {
		if err != nil {
			handleError(w, http.StatusBadRequest, "could not get apikey", err.Error())
			return
		}
	}
	apiKey := fmt.Sprintf("%v", _apiKey)

	if err := session.UpdateAPIKey(apiKey, accountID); err != nil {
		handleError(w, http.StatusUnauthorized, "could not update apikey", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

func deleteAPITokenHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_accountID, ok := body["accountID"]
	if !ok {
		if err != nil {
			handleError(w, http.StatusBadRequest, "could not get account id", err.Error())
			return
		}
	}
	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.DeleteAPIKey(accountID); err != nil {
		handleError(w, http.StatusUnauthorized, "could not delete apikey", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

func getScheduleHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
	log.Println("here")

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}

	docs, err := session.GetDocument(accountID)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get docs", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(docs)
	return
}

func setScheduleHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}

	if err := session.CreateDocument(accountID, body); err != nil {
		handleError(w, http.StatusUnauthorized, "could not create record", err.Error())
		return
	}
}

// TODO: complete this
func deleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
}

func updateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
}

func getCloudSessions(r *http.Request) (*ibmcloud.Session, error) {
	accessTokenVal, err := r.Cookie(accessToken)
	if err != nil {
		return nil, err
	}
	refreshTokenVal, err := r.Cookie(refreshToken)
	if err != nil {
		return nil, err
	}
	expirationValStr, err := r.Cookie(expiration)
	if err != nil {
		return nil, err
	}

	expirationVal, err := strconv.Atoi(expirationValStr.Value)
	if err != nil {
		return nil, err
	}

	session := &ibmcloud.Session{
		Token: &ibmcloud.Token{
			AccessToken:  accessTokenVal.Value,
			RefreshToken: refreshTokenVal.Value,
			Expiration:   expirationVal,
		},
	}

	return session.RenewSession()
}

func handleError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, fmt.Sprintf(errorMessageFormat, strings.Join(message, " ")))
}
