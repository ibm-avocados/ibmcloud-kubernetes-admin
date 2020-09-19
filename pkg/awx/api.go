package awx

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