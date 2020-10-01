package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/awx"
)

func GetAWXWorkflowJobTemplates(c echo.Context) error {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	templates, err := awx.GetWorkflowJobTemplates(token)
	if err != nil {
		return err
	}

	query := c.QueryParam("labels")
	fmt.Println("QUERY: ", query)

	var result []awx.ResultsWorkflowTemplate
	if query == "" {
		result = templates.Results
	} else {
		for _, res := range templates.Results {
			for _, label := range res.SummaryFields.Labels.Results {
				if label.Name == query {
					result = append(result, res)
					break
				}
			}
		}
	}

	return c.JSON(http.StatusOK, result)
}

func GetAWXJobTemplates(c echo.Context) error {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	templates, err := awx.GetJobTemplates(token)
	if err != nil {
		return err
	}

	result := templates.Results

	return c.JSON(http.StatusOK, result)
}

func LaunchAWXWorkflowJobTemplate(c echo.Context) error {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	var body awx.WorkflowJobTeplatesLaunchBody
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(&body)
	if err != nil {
		return err
	}

	res, err := awx.LaunchWorkflowJobTemplate(token, body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
