package kubeadmin

import (
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/vcs"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/vcs/github"
	"github.com/labstack/echo/v4"
)

func GithubCommentHandler(c echo.Context) error {
	comment := new(vcs.GithubIssueComment)
	c.Request().Header.Add("Content-Type", "application/json")

	if err := c.Bind(comment); err != nil {
		return err
	}

	err := github.CreateComment(*comment)

	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "OK")
}
