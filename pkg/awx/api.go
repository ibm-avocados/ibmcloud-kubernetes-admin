package awx

import "encoding/json"

const (
	protocol    = "https://"
	awxEndpoint = "awx.ibmdeveloper.net"
)

func getWorkflowJobTemplates(token string) (*WorkflowJobTemplates, error) {
	var res WorkflowJobTemplates
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	endpoint := protocol + awxEndpoint + "/api/v2/workflow_job_templates/"
	fetch(endpoint, header, nil, &res)
	return &res, nil
}

func getJobTemplates(token string) (*JobTemplates, error) {
	var res JobTemplates
	header := map[string]string{
		"Authorization": "Bearer " + token,
	}
	endpoint := protocol + awxEndpoint + "/api/v2/job_templates/"
	fetch(endpoint, header, nil, &res)
	return &res, nil
}

func launchWorkflowJobTemplate(token string, body WorkflowJobTeplatesLaunchBody) (interface{}, error) {
	var res interface{}
	header := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	endpont := protocol + awxEndpoint + "/api/v2/workflow_job_templates/" + body.ID + "/launch/"

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	err = postBody(endpont, header, nil, b, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
