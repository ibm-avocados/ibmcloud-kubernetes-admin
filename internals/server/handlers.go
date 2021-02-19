package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

const (
	errorMessageFormat = `{"msg": "error: %s"}`
	statusOkMessage    = `{"status": "ok"}`
	sessionName        = "cloud_session"
	accessTokenKey     = "access_token"
	refreshTokenKey    = "refresh_token"
	expiration         = "expiration"
	cookiePath         = "/api/v1"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	handleError(w, http.StatusNotFound, "not found")
}

func setCookie(c echo.Context, session *ibmcloud.Session) {
	accessTokenCookie := &http.Cookie{Name: accessTokenKey, Value: session.Token.AccessToken, Path: cookiePath}
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := &http.Cookie{Name: refreshTokenKey, Value: session.Token.RefreshToken, Path: cookiePath}
	c.SetCookie(refreshTokenCookie)

	expirationStr := strconv.Itoa(session.Token.Expiration)

	expirationCookie := &http.Cookie{Name: expiration, Value: expirationStr, Path: cookiePath}
	c.SetCookie(expirationCookie)
}

func deleteCookie(c echo.Context) {
	accessTokenCookie := &http.Cookie{Name: accessTokenKey, Value: "", Path: cookiePath}
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := &http.Cookie{Name: refreshTokenKey, Value: "", Path: cookiePath}
	c.SetCookie(refreshTokenCookie)

	expirationCookie := &http.Cookie{Name: expiration, Value: "", Path: cookiePath}
	c.SetCookie(expirationCookie)
}

func getCloudSessions(c echo.Context) (*ibmcloud.Session, error) {
	var accessToken string
	var refreshToken string
	var expirationTime int
	accessTokenVal, err := c.Cookie(accessTokenKey)
	if err != nil {
		bearerToken := c.Request().Header.Get("Authorization")
		if bearerToken == "" {
			return nil, err
		}
		parsedToken := strings.Split(bearerToken, " ")
		if len(parsedToken) != 2 {
			return nil, err
		}
		accessToken = parsedToken[1]
	} else {
		accessToken = accessTokenVal.Value
	}

	refreshTokenVal, err := c.Cookie(refreshTokenKey)
	if err != nil {
		refreshToken = c.Request().Header.Get("X-Auth-Refresh-Token")
		if refreshToken == "" {
			return nil, err
		}
	} else {
		refreshToken = refreshTokenVal.Value
	}

	expirationValStr, err := c.Cookie(expiration)
	if err != nil {
		expiration := c.Request().Header.Get("Expiration")
		if expiration == "" {
			expirationTime = 0
		}
	} else {
		expirationTime, err = strconv.Atoi(expirationValStr.Value)
		if err != nil {
			return nil, err
		}
	}

	session := &ibmcloud.Session{
		Token: &ibmcloud.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Expiration:   expirationTime,
		},
	}

	return session.RenewSession()
}

func handleError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, fmt.Sprintf(errorMessageFormat, strings.Join(message, " ")))
}
