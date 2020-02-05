package main

import (
	"encoding/gob"
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

func clusterListHandler(w http.ResponseWriter, r *http.Request) {
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

	clusters, err := session.GetClusters(accountID, "")
	if err != nil {
		handleError(w, http.StatusUnauthorized, "could not get clusters", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(clusters)
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
