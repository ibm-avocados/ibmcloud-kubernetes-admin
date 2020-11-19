package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/github"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func GithubCommentHandler(c echo.Context) error {
	comment := new(ibmcloud.GithubIssueComment)
	c.Request().Header.Add("Content-Type", "application/json")

	if err := c.Bind(comment); err != nil {
		return err
	}

	err := github.CreateComment(*comment, "templates/message.gotmpl")

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, statusOkMessage)
}
