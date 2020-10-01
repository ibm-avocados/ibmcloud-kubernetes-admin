package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateAdminEmails(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}
	var body AccountEmailBody
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	fmt.Println(body)
	if err := session.CreateAdminEmails(body.AccountID, body.Email...); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func AddAdminEmails(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}
	var body AccountEmailBody
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	if err := session.AddAdminEmails(body.AccountID, body.Email...); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func RemoveAdminEmails(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}
	var body AccountEmailBody
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	if err := session.RemoveAdminEmails(body.AccountID, body.Email...); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func DeleteAdminEmails(c echo.Context) error {
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

	_accountID, ok := body["accountID"]
	if !ok {
		return errors.New("no valid account ID")
	}

	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.DeleteAdminEmails(accountID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func GetAdminEmails(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	emails, err := session.GetAccountAdminEmails(accountID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, emails)
}
