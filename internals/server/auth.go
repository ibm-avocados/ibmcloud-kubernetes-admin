package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func AuthenticationHandler(c echo.Context) error {
	accountLogin := new(AccountLogin)
	if err := c.Bind(accountLogin); err != nil {
		fmt.Println("1", err)
		return err
	}

	session, err := ibmcloud.Authenticate(accountLogin.OTP)
	if err != nil {
		fmt.Println("2", err)
		return err
	}

	setCookie(c, session)

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func AuthenticationWithAccountHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	if !session.IsValid() {
		return err
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	accountID := fmt.Sprintf("%v", body["id"])

	log.Println("Account id", accountID)

	accountSession, err := session.BindAccountToToken(accountID)
	if err != nil {
		return err
	}

	setCookie(c, accountSession)
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func LoginHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	if !session.IsValid() {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func TokenEndpointHandler(c echo.Context) error {
	endpoints, err := ibmcloud.GetIdentityEndpoints()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, endpoints)
}
