package kubeadmin

import (
	"net/http"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vcs"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vcs/github"
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
