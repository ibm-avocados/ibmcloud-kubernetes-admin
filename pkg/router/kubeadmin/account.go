package kubeadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func ResourceGroupHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

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
}

func AccountListHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

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
}

func CheckApiKeyHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

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
}
