package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessGroupsHandler(c echo.Context) error {
	session, err := getCloudSessions(c)
	if err != nil {
		return err
	}

	accountID := c.Param("accountID")

	accessGroups, err := session.GetAccessGroups(accountID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accessGroups)
}
