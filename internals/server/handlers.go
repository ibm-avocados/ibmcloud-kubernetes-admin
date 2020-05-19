package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
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

func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	handleError(w, http.StatusNotFound, "not found")
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

func handleError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, fmt.Sprintf(errorMessageFormat, strings.Join(message, " ")))
}
