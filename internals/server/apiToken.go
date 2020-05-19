package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) SetAPITokenHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) CheckAPITokenHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) UpdateAPITokenHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) DeleteAPITokenHandler(w http.ResponseWriter, r *http.Request) {
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
