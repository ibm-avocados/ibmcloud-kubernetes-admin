package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func DeleteTagHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(c.Request().Body)

	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	res, err := session.DeleteTag(body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: fmt.Sprintf("success. %d tags deleted", len(res.Results))})
}

func GetTagHandler(c echo.Context) error {
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

	crn, ok := body["crn"]
	if !ok {
		return errors.New("crn not in body")
	}

	clusterCRN := fmt.Sprintf("%v", crn)

	tags, err := session.GetTags(clusterCRN)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tags)
}

func SetClusterTagHandler(c echo.Context) error {
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

	tag, ok := body["tag"]
	if !ok {
		return errors.New("tag not in body")
	}

	resourceGroup, ok := body["resourceGroup"]
	if !ok {
		return errors.New("resourceGroup not in body")
	}

	clusterID := c.Param("clusterID")

	clusterTag := fmt.Sprintf("%v", tag)
	clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)

	res, err := session.SetClusterTag(clusterTag, clusterID, clusterResourceGroup)
	if err != nil {
		log.Println("set cluster tag : ", err)
		return err
	}

	return c.JSON(http.StatusOK, StatusOK{Message: fmt.Sprintf("success, wrote %d tags", len(res.Results))})
}

func SetTagHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body ibmcloud.UpdateTag
	decoder := json.NewDecoder(c.Request().Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	res, err := session.SetTag(body)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: fmt.Sprintf("success, wrote %d tags", len(res.Results))})
}
