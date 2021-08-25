package kubeadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func DeleteTagHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(c.Request().Body)

		if err != nil {
			return err
		}

		res, err := session.DeleteTag(bytes)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, StatusOK{Message: fmt.Sprintf("success. %d tags deleted", len(res.Results))})
	}
}
func GetTagHandler(provider infra.SessionProvider) echo.HandlerFunc {
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
}
func SetClusterTagHandler(provider infra.SessionProvider) echo.HandlerFunc {
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
}
func SetTagHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		res, err := session.SetTag(bytes)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, StatusOK{Message: fmt.Sprintf("success, wrote %d tags", len(res.Results))})
	}
}
