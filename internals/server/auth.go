package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func (s *Server) AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
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

	setCookie(w, session)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) AuthenticationWithAccountHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	if !session.IsValid() {
		handleError(w, http.StatusUnauthorized, "session expired")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) TokenEndpointHandler(w http.ResponseWriter, r *http.Request) {
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
