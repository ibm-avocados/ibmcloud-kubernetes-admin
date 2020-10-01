package server

import (
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
