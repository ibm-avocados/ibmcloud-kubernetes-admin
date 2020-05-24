package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// func (s *Server) GetScheduleHandler(w http.ResponseWriter, r *http.Request) {
// 	session, err := getCloudSessions(r)
// 	if err != nil {
// 		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
// 		return
// 	}

// 	vars := mux.Vars(r)

// 	accountID, ok := vars["accountID"]

// 	if !ok {
// 		handleError(w, http.StatusBadRequest, "could not get accountID")
// 		return
// 	}

// 	docs, err := session.GetDocument(accountID)
// 	if err != nil {
// 		handleError(w, http.StatusUnauthorized, "could not get docs", err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	e := json.NewEncoder(w)
// 	e.Encode(docs)
// 	return
// }

func (s *Server) SetScheduleHandler(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, statusOkMessage)
}

// TODO: complete this
func (s *Server) DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
}

func (s *Server) UpdateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get session", err.Error())
		return
	}
}

func (s *Server) GetAllScheduleHandler(w http.ResponseWriter, r *http.Request) {
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

	docs, err := session.GetAllDocument(accountID)
	if err != nil {
		handleError(w, http.StatusBadRequest, "could not get accountID")
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(docs)
	return
}
