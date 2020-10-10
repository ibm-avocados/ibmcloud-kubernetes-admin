package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ResourceGroupHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	accountResources, err := session.GetAccountResources(accountID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accountResources)
}

func AccountListHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		log.Printf("could not get session %v\n", err)
		return err
	}

	accounts, err := session.GetAccounts()
	if err != nil {
		log.Printf("could not get accounts using access token %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}

func CheckApiKeyHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	_apikey, ok := body["apikey"]
	if !ok {
		return errors.New("no apikey attached to body")
	}

	apikey := fmt.Sprintf("%v", _apikey)

	apikeyDetails, err := session.CheckToken(apikey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, apikeyDetails)
}
