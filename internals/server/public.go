package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func VersionEndpointHandler(c echo.Context) error {
	versions, err := ibmcloud.GetVersions()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, versions)
}

func LocationEndpointHandler(c echo.Context) error {
	locations, err := ibmcloud.GetLocations()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, locations)
}

func LocationGeoEndpointHandler(c echo.Context) error {
	geo := c.Param("geo")
	locations, err := ibmcloud.GetGeoLocations(geo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, locations)
}

func ZonesEndpointHandler(c echo.Context) error {
	showFlavors := c.QueryParam("showFlavors")
	location := c.QueryParam("location")
	zones, err := ibmcloud.GetZones(showFlavors, location)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, zones)
}

func MachineTypeHandler(c echo.Context) error {
	serverType := c.QueryParam("type")
	os := c.QueryParam("os")
	cpuLimitStr := c.QueryParam("cpuLimit")
	cpuLimit, err := strconv.Atoi(cpuLimitStr)
	if err != nil {
		return err
	}
	memoryLimitStr := c.QueryParam("memoryLimit")
	memoryLimit, err := strconv.Atoi(memoryLimitStr)
	if err != nil {
		return err
	}

	datacenter := c.Param("datacenter")

	flavors, err := ibmcloud.GetMachineType(datacenter, serverType, os, cpuLimit, memoryLimit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, flavors)
}
