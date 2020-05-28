package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetMetaDataHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}

	metadata, err := session.GetAccountMetaData(accountID)
	log.Println(metadata)
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(metadata)
}

func (s *Server) UpdateMetaDataHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_org, ok := body["org"]
	if !ok {
		handleError(w, http.StatusBadRequest, "org not found", err.Error())
		return
	}
	org := fmt.Sprintf("%v", _org)

	_space, ok := body["space"]
	if !ok {
		handleError(w, http.StatusBadRequest, "space not found", err.Error())
		return
	}

	space := fmt.Sprintf("%v", _space)

	_region, ok := body["region"]
	if !ok {
		handleError(w, http.StatusBadRequest, "region not found", err.Error())
		return
	}

	region := fmt.Sprintf("%v", _region)

	_accessGroup, ok := body["accessGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "accessGroup not found", err.Error())
		return
	}

	accessGroup := fmt.Sprintf("%v", _accessGroup)

	_issueRepo, ok := body["issueRepo"]
	if !ok {
		handleError(w, http.StatusBadRequest, "issueRepo not found", err.Error())
		return
	}

	issueRepo := fmt.Sprintf("%v", _issueRepo)

	_grantClusterRepo, ok := body["grantClusterRepo"]
	if !ok {
		handleError(w, http.StatusBadRequest, "grantClusterRepo not found", err.Error())
		return
	}

	grantClusterRepo := fmt.Sprintf("%v", _grantClusterRepo)

	_githubUser, ok := body["githubUser"]
	if !ok {
		handleError(w, http.StatusBadRequest, "githubUser not found", err.Error())
		return
	}

	githubUser := fmt.Sprintf("%v", _githubUser)

	_githubToken, ok := body["githubToken"]
	if !ok {
		handleError(w, http.StatusBadRequest, "githubToken not found", err.Error())
		return
	}

	githubToken := fmt.Sprintf("%v", _githubToken)

	if err := session.UpdateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken); err != nil {
		handleError(w, http.StatusUnauthorized, "could not create record", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

func (s *Server) CreateMetaDataHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]

	if !ok {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not decode", err.Error())
		return
	}

	_org, ok := body["org"]
	if !ok {
		handleError(w, http.StatusBadRequest, "org not found", err.Error())
		return
	}
	org := fmt.Sprintf("%v", _org)

	_space, ok := body["space"]
	if !ok {
		handleError(w, http.StatusBadRequest, "space not found", err.Error())
		return
	}

	space := fmt.Sprintf("%v", _space)

	_region, ok := body["region"]
	if !ok {
		handleError(w, http.StatusBadRequest, "region not found", err.Error())
		return
	}

	region := fmt.Sprintf("%v", _region)

	_accessGroup, ok := body["accessGroup"]
	if !ok {
		handleError(w, http.StatusBadRequest, "accessGroup not found", err.Error())
		return
	}

	accessGroup := fmt.Sprintf("%v", _accessGroup)

	_issueRepo, ok := body["issueRepo"]
	if !ok {
		handleError(w, http.StatusBadRequest, "issueRepo not found", err.Error())
		return
	}

	issueRepo := fmt.Sprintf("%v", _issueRepo)

	_grantClusterRepo, ok := body["grantClusterRepo"]
	if !ok {
		handleError(w, http.StatusBadRequest, "grantClusterRepo not found", err.Error())
		return
	}

	grantClusterRepo := fmt.Sprintf("%v", _grantClusterRepo)

	_githubUser, ok := body["githubUser"]
	if !ok {
		handleError(w, http.StatusBadRequest, "githubUser not found", err.Error())
		return
	}

	githubUser := fmt.Sprintf("%v", _githubUser)

	_githubToken, ok := body["githubToken"]
	if !ok {
		handleError(w, http.StatusBadRequest, "githubToken not found", err.Error())
		return
	}

	githubToken := fmt.Sprintf("%v", _githubToken)

	if err := session.CreateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken); err != nil {
		handleError(w, http.StatusUnauthorized, "could not create record", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}
