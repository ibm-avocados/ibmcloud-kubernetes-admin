package kubeadmin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/awx"
	"github.com/labstack/echo/v4"
)

func GetAWXWorkflowJobTemplates(c echo.Context) error {
	token := os.Getenv("AWX_ACCESS_TOKEN")
	templates, err := awx.GetWorkflowJobTemplates(token)
	if err != nil {
		return err
	}

	query := c.QueryParam("labels")
	fmt.Println("QUERY: ", query)

	result := make([]awx.ResultsWorkflowTemplate, 0)
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
		fmt.Println("1:", err)
		return err
	}

	res, err := awx.LaunchWorkflowJobTemplate(token, body)
	if err != nil {
		fmt.Println("2:", err)

		return err
	}

	return c.JSON(http.StatusOK, res)
}

func GetGrantClusterID(c echo.Context) error {
	grantclusterWorkflowID, ok := os.LookupEnv("GRANT_CLUSTER_WORKFLOW_ID")
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "grant cluster workflow id not set")
	}
	data := struct {
		ID string `json:"id"`
	}{
		ID: grantclusterWorkflowID,
	}
	return c.JSON(http.StatusOK, data)
}
