package awx

func GetWorkflowJobTemplates(token string) (*WorkflowJobTemplates, error) {
	return getWorkflowJobTemplates(token)
}

func GetJobTemplates(token string) (*JobTemplates, error) {
	return getJobTemplates(token)
}
