package kubeadmin

import (
	"net/http"
	"strconv"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func VersionEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		versions, err := session.GetVersions()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, versions)
	}
}

func LocationEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		locations, err := session.GetLocations()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, locations)
	}
}

func LocationGeoEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		geo := c.Param("geo")
		locations, err := session.GetGeoLocations(geo)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, locations)
	}
}
func ZonesEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		showFlavors := c.QueryParam("showFlavors")
		location := c.QueryParam("location")

		zones, err := session.GetZones(showFlavors, location)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, zones)
	}
}

func MachineTypeHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

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

		flavors, err := session.GetMachineType(datacenter, serverType, os, cpuLimit, memoryLimit)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, flavors)
	}
}
