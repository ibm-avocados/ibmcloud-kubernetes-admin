package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetAPITokenHandler(c echo.Context) error {
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
		if err != nil {
			return err
		}
	}
	accountID := fmt.Sprintf("%v", _accountID)

	_apiKey, ok := body["apiKey"]
	if !ok {
		if err != nil {
			return err
		}
	}
	apiKey := fmt.Sprintf("%v", _apiKey)

	if err := session.SetAPIKey(apiKey, accountID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func CheckAPITokenHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&body); err != nil {
		return err
	}

	_accountID, ok := body["accountID"]
	if !ok {
		return err
	}
	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.CheckAPIKey(accountID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func UpdateAPITokenHandler(c echo.Context) error {
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
		return errors.New("No valid account ID")
	}
	accountID := fmt.Sprintf("%v", _accountID)

	_apiKey, ok := body["apiKey"]
	if !ok {
		return errors.New("No valid apikey ID")
	}
	apiKey := fmt.Sprintf("%v", _apiKey)

	if err := session.UpdateAPIKey(apiKey, accountID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func DeleteAPITokenHandler(c echo.Context) error {
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
		return errors.New("No valid account ID")
	}
	accountID := fmt.Sprintf("%v", _accountID)

	if err := session.DeleteAPIKey(accountID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}
