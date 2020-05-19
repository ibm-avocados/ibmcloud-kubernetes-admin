package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) ResourceGroupHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) AccountListHandler(w http.ResponseWriter, r *http.Request) {
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
