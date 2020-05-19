package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) GetBillingHandler(w http.ResponseWriter, r *http.Request) {
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
