package kubeadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func ClusterDeleteHandler(provider infra.SessionProvider) echo.HandlerFunc {
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

		id := fmt.Sprintf("%v", body["id"])
		resoueceGroup := fmt.Sprintf("%v", body["resourceGroup"])
		deleteResources := fmt.Sprintf("%v", body["deleteResources"])
		err = session.DeleteCluster(id, resoueceGroup, deleteResources)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, StatusOK{Message: "success"})
	}
}

func ClusterCreateHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}

		createResponse, err := session.CreateCluster(bytes)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, createResponse)
	}
}

func LocationGeoEndpointInfoHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		clusters, err := session.GetClusters("")
		if err != nil {
			return err
		}

		res := make(map[string]int)

		for _, c := range clusters {
			if count, ok := res[c.DataCenter]; ok {
				res[c.DataCenter] = count + 1
			} else {
				res[c.DataCenter] = 1
			}
		}
		return c.JSON(http.StatusOK, res)
	}
}

func ClusterListHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		clusters, err := session.GetClusters("")
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, clusters)
	}
}
func ClusterHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		clusterID := c.Param("clusterID")

		var body map[string]interface{}
		decoder := json.NewDecoder(c.Request().Body)
		err = decoder.Decode(&body)
		if err != nil {
			return err
		}

		resourceGroup, ok := body["resourceGroup"]
		if !ok {
			return errors.New("no resourceGroup attached to body")
		}
		clusterResourceGroup := fmt.Sprintf("%v", resourceGroup)
		cluster, err := session.GetCluster(clusterID, clusterResourceGroup)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, cluster)
	}
}

func ClusterWorkerListHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		clusterID := c.Param("clusterID")

		workers, err := session.GetWorkers(clusterID)

		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, workers)
	}
}

func VlanEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		datacenter := c.Param("datacenter")

		vlans, err := session.GetDatacenterVlan(datacenter)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, vlans)
	}
}
