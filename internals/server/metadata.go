package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetMetaDataHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	metadata, err := session.GetAccountMetaData(accountID)
	return c.JSON(http.StatusOK, metadata)
}

func UpdateMetaDataHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	_org, ok := body["org"]
	if !ok {
		return errors.New("org not found in body")
	}
	org := fmt.Sprintf("%v", _org)

	_space, ok := body["space"]
	if !ok {
		return errors.New("space not found in body")
	}

	space := fmt.Sprintf("%v", _space)

	_region, ok := body["region"]
	if !ok {
		return errors.New("region not found in body")
	}

	region := fmt.Sprintf("%v", _region)

	_accessGroup, ok := body["accessGroup"]
	if !ok {
		return errors.New("accessGroup not found in body")
	}

	accessGroup := fmt.Sprintf("%v", _accessGroup)

	_issueRepo, ok := body["issueRepo"]
	if !ok {
		return errors.New("issueRepo not found in body")
	}

	issueRepo := fmt.Sprintf("%v", _issueRepo)

	_grantClusterRepo, ok := body["grantClusterRepo"]
	if !ok {
		return errors.New("grantClusterRepo not found in body")
	}

	grantClusterRepo := fmt.Sprintf("%v", _grantClusterRepo)

	_githubUser, ok := body["githubUser"]
	if !ok {
		return errors.New("githubUser not found in body")
	}

	githubUser := fmt.Sprintf("%v", _githubUser)

	_githubToken, ok := body["githubToken"]
	if !ok {
		return errors.New("githubToken not found in body")
	}

	githubToken := fmt.Sprintf("%v", _githubToken)

	if err := session.UpdateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func CreateMetaDataHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	var body map[string]interface{}
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	_org, ok := body["org"]
	if !ok {
		return errors.New("org not found in body")
	}
	org := fmt.Sprintf("%v", _org)

	_space, ok := body["space"]
	if !ok {
		return errors.New("space not found in body")
	}

	space := fmt.Sprintf("%v", _space)

	_region, ok := body["region"]
	if !ok {
		return errors.New("region not found in body")
	}

	region := fmt.Sprintf("%v", _region)

	_accessGroup, ok := body["accessGroup"]
	if !ok {
		return errors.New("accessGroup not found in body")
	}

	accessGroup := fmt.Sprintf("%v", _accessGroup)

	_issueRepo, ok := body["issueRepo"]
	if !ok {
		return errors.New("issueRepo not found in body")
	}

	issueRepo := fmt.Sprintf("%v", _issueRepo)

	_grantClusterRepo, ok := body["grantClusterRepo"]
	if !ok {
		return errors.New("grantClusterRepo not found in body")
	}

	grantClusterRepo := fmt.Sprintf("%v", _grantClusterRepo)

	_githubUser, ok := body["githubUser"]
	if !ok {
		return errors.New("githubUser not found in body")
	}

	githubUser := fmt.Sprintf("%v", _githubUser)

	_githubToken, ok := body["githubToken"]
	if !ok {
		return errors.New("githubToken not found in body")
	}

	githubToken := fmt.Sprintf("%v", _githubToken)

	if err := session.CreateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}
