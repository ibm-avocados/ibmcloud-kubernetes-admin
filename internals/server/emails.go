package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) CreateAdminEmails(w http.ResponseWriter, r *http.Request) {
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

	_accountID, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no accountID provided", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", _accountID)

	_email, ok := body["email"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no emails provided", err.Error())
		return
	}

	emails, ok := _email.([]string)
	if !ok {
		handleError(w, http.StatusBadRequest, "not a valid type", err.Error())
		return
	}

	if err := session.CreateAdminEmails(accountID, emails...); err != nil {
		handleError(w, http.StatusInternalServerError, "could not create", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) AddAdminEmails(w http.ResponseWriter, r *http.Request) {
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

	_accountID, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no accountID provided", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", _accountID)

	_email, ok := body["email"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no emails provided", err.Error())
		return
	}

	emails, ok := _email.([]string)
	if !ok {
		handleError(w, http.StatusBadRequest, "not a valid type", err.Error())
		return
	}

	if err := session.AddAdminEmails(accountID, emails...); err != nil {
		handleError(w, http.StatusInternalServerError, "could not add", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) RemoveAdminEmails(w http.ResponseWriter, r *http.Request) {
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

	_accountID, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no accountID provided", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", _accountID)

	_email, ok := body["email"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no emails provided", err.Error())
		return
	}

	emails, ok := _email.([]string)
	if !ok {
		handleError(w, http.StatusBadRequest, "not a valid type", err.Error())
		return
	}

	if err := session.RemoveAdminEmails(accountID, emails...); err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) DeleteAdminEmails(w http.ResponseWriter, r *http.Request) {
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

	_accountID, ok := body["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no accountID provided", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.DeleteAdminEmails(accountID); err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func (s *Server) GetAdminEmails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}

	vars := mux.Vars(r)

	accountID, ok := vars["accountID"]
	if !ok {
		handleError(w, http.StatusBadRequest, "no accountID provided", err.Error())
		return
	}

	emails, err := session.GetAccountAdminEmails(accountID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not delete", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(emails)
}
