package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/moficodes/ibmcloud-kubernetes-admin/ibmcloud"
)

const (
	errorMessageFormat = `{"msg": "error: %s"}`
	statusOkMessage    = `{"status": "ok"}`
	sessionID          = "ibmcloud_token"
	sessionName        = "cloud_session"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func init() {
	gob.Register(&ibmcloud.Session{})
}

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

func authenticationWithAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	session, err := store.Get(r, sessionID)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get session", err.Error())
		return
	}

	cloudSession, err := getCloudSessions(r, session)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}

	if !cloudSession.IsValid() {
		handleError(w, http.StatusUnauthorized, "session not valid")
		return
	}

	session.Options = &sessions.Options{
		MaxAge: 3600 * 60 * 24,
	}

	var account ibmcloud.Account
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&account)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not decode", err.Error())
		return
	}

	accountID := account.Metadata.GUID

	fmt.Println(accountID)

	accountSession, err := cloudSession.BindAccountToToken(account)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not bind account to token", err.Error())
		return
	}

	session.Values[sessionName] = accountSession

	err = sessions.Save(r, w)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not save session", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	session, err := store.Get(r, sessionID)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get session", err.Error())
		return
	}

	session.Options = &sessions.Options{
		MaxAge: 3600 * 60 * 24,
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not decode", err.Error())
		return
	}

	otp := fmt.Sprintf("%v", body["otp"])

	fmt.Println(otp)
	cloudSession, err := ibmcloud.Authenticate(otp)
	if err != nil {
		log.Println("could not authenticate with the otp provided")
		log.Println(err.Error())
		handleError(w, http.StatusInternalServerError, "could not authenticate with the otp provided", err.Error())
		return
	}

	session.Values[sessionName] = cloudSession
	err = sessions.Save(r, w)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not save session", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	session, err := store.Get(r, sessionID)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get session", err.Error())
		return
	}
	cloudSession, err := getCloudSessions(r, session)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}
	fmt.Println(cloudSession.Token.Expiration)

	if !cloudSession.IsValid() {
		handleError(w, http.StatusUnauthorized, "session expired")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func accountListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionID)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get session", err.Error())
		return
	}
	cloudSession, err := getCloudSessions(r, session)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}

	accounts, err := cloudSession.GetAccounts()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get accounts using access token", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(accounts)
}

func handleError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, fmt.Sprintf(errorMessageFormat, strings.Join(message, " ")))
}

func getCloudSessions(r *http.Request, session *sessions.Session) (*ibmcloud.Session, error) {
	val := session.Values[sessionName]
	var cloudSession *ibmcloud.Session
	var ok bool
	if cloudSession, ok = val.(*ibmcloud.Session); !ok {
		return nil, errors.New("could not cast session to cloud session object")
	}
	return cloudSession, nil
}
