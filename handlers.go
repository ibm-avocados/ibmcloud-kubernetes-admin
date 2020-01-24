package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
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
		handleError(w, http.StatusInternalServerError, "could not decode", err.Error())
		return
	}

	accountID := fmt.Sprintf("%v", body["id"])

	fmt.Println("Account id", accountID)

	accountSession, err := session.BindAccountToToken(accountID)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not bind account to token", err.Error())
		return
	}

	setCookie(w, accountSession)

	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not save session", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusOkMessage)
}

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not decode", err.Error())
		return
	}

	otp := fmt.Sprintf("%v", body["otp"])

	fmt.Println(otp)
	session, err := ibmcloud.Authenticate(otp)
	if err != nil {
		log.Println("could not authenticate with the otp provided")
		log.Println(err.Error())
		handleError(w, http.StatusInternalServerError, "could not authenticate with the otp provided", err.Error())
		return
	}

	fmt.Println(session.Token.Expiration)

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
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}
	fmt.Println(session.Token.Expiration)

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
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}

	accounts, err := session.GetAccounts()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "could not get accounts using access token", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(accounts)
}

func clusterListHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getCloudSessions(r)
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get session", err.Error())
		return
	}
	clusters, err := session.GetClusters("")
	if err != nil {
		handleError(w, http.StatusNotFound, "could not get clusters", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(clusters)
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, statusOkMessage)
}

func handleError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, fmt.Sprintf(errorMessageFormat, strings.Join(message, " ")))
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
	expirationVal, _ := strconv.Atoi(expirationValStr.Value)
	session := &ibmcloud.Session{
		Token: &ibmcloud.Token{
			AccessToken:  accessTokenVal.Value,
			RefreshToken: refreshTokenVal.Value,
			Expiration:   expirationVal,
		},
	}

	return session, nil
}
