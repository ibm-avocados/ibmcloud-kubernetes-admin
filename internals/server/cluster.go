package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func ClusterDeleteHandler(c echo.Context) error {
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

	id := fmt.Sprintf("%v", body["id"])
	resoueceGroup := fmt.Sprintf("%v", body["resourceGroup"])
	deleteResources := fmt.Sprintf("%v", body["deleteResources"])
	err = session.DeleteCluster(id, resoueceGroup, deleteResources)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, StatusOK{Message: "success"})
}

func ClusterCreateHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	var body ibmcloud.CreateClusterRequest

	decoder := json.NewDecoder(c.Request().Body)

	err = decoder.Decode(&body)
	if err != nil {
		return err
	}

	createResponse, err := session.CreateCluster(body)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createResponse)
}

func LocationGeoEndpointInfoHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
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

func ClusterListHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	clusters, err := session.GetClusters("")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, clusters)
}

func ClusterHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
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

func ClusterWorkerListHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
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

func VlanEndpointHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
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
