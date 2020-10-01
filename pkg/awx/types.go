package awx

import "time"

// WorkflowJobTemplate represents a workflow job template
type WorkflowJobTemplates struct {
	Count    int                       `json:"count"`
	Next     interface{}               `json:"next"`
	Previous interface{}               `json:"previous"`
	Results  []ResultsWorkflowTemplate `json:"results"`
}
type RelatedWorkflowTemplates struct {
	CreatedBy                      string `json:"created_by"`
	ModifiedBy                     string `json:"modified_by"`
	LastJob                        string `json:"last_job"`
	WorkflowJobs                   string `json:"workflow_jobs"`
	Schedules                      string `json:"schedules"`
	Launch                         string `json:"launch"`
	WebhookKey                     string `json:"webhook_key"`
	WebhookReceiver                string `json:"webhook_receiver"`
	WorkflowNodes                  string `json:"workflow_nodes"`
	Labels                         string `json:"labels"`
	ActivityStream                 string `json:"activity_stream"`
	NotificationTemplatesStarted   string `json:"notification_templates_started"`
	NotificationTemplatesSuccess   string `json:"notification_templates_success"`
	NotificationTemplatesError     string `json:"notification_templates_error"`
	NotificationTemplatesApprovals string `json:"notification_templates_approvals"`
	AccessList                     string `json:"access_list"`
	ObjectRoles                    string `json:"object_roles"`
	SurveySpec                     string `json:"survey_spec"`
	Copy                           string `json:"copy"`
}
type LastJob struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Finished    time.Time `json:"finished"`
	Status      string    `json:"status"`
	Failed      bool      `json:"failed"`
}
type LastUpdate struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Failed      bool   `json:"failed"`
}
type CreatedBy struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type ModifiedBy struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ExecuteRole struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}
type ReadRole struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}
type ApprovalRole struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}
type ObjectRoles struct {
	AdminRole    AdminRole    `json:"admin_role"`
	ExecuteRole  ExecuteRole  `json:"execute_role"`
	ReadRole     ReadRole     `json:"read_role"`
	ApprovalRole ApprovalRole `json:"approval_role"`
}

type Labels struct {
	Count   int            `json:"count"`
	Results []LabelResults `json:"results"`
}

type LabelResults struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RecentJobs struct {
	ID         int         `json:"id"`
	Status     string      `json:"status"`
	Finished   time.Time   `json:"finished"`
	CanceledOn interface{} `json:"canceled_on"`
	Type       string      `json:"type"`
}
type Inventory struct {
	ID                           int    `json:"id"`
	Name                         string `json:"name"`
	Description                  string `json:"description"`
	HasActiveFailures            bool   `json:"has_active_failures"`
	TotalHosts                   int    `json:"total_hosts"`
	HostsWithActiveFailures      int    `json:"hosts_with_active_failures"`
	TotalGroups                  int    `json:"total_groups"`
	HasInventorySources          bool   `json:"has_inventory_sources"`
	TotalInventorySources        int    `json:"total_inventory_sources"`
	InventorySourcesWithFailures int    `json:"inventory_sources_with_failures"`
	OrganizationID               int    `json:"organization_id"`
	Kind                         string `json:"kind"`
}

type ResultsWorkflowTemplate struct {
	ID                   int           `json:"id"`
	Type                 string        `json:"type"`
	URL                  string        `json:"url"`
	Related              Related       `json:"related"`
	Created              time.Time     `json:"created"`
	Modified             time.Time     `json:"modified"`
	Name                 string        `json:"name"`
	Description          string        `json:"description"`
	LastJobRun           time.Time     `json:"last_job_run"`
	LastJobFailed        bool          `json:"last_job_failed"`
	NextJobRun           interface{}   `json:"next_job_run"`
	Status               string        `json:"status"`
	ExtraVars            string        `json:"extra_vars"`
	Organization         interface{}   `json:"organization"`
	SurveyEnabled        bool          `json:"survey_enabled"`
	AllowSimultaneous    bool          `json:"allow_simultaneous"`
	AskVariablesOnLaunch bool          `json:"ask_variables_on_launch"`
	Inventory            interface{}   `json:"inventory"`
	Limit                interface{}   `json:"limit"`
	ScmBranch            interface{}   `json:"scm_branch"`
	AskInventoryOnLaunch bool          `json:"ask_inventory_on_launch"`
	AskScmBranchOnLaunch bool          `json:"ask_scm_branch_on_launch"`
	AskLimitOnLaunch     bool          `json:"ask_limit_on_launch"`
	WebhookService       string        `json:"webhook_service"`
	WebhookCredential    interface{}   `json:"webhook_credential"`
	SummaryFields        SummaryFields `json:"summary_fields,omitempty"`
}

type JobTemplates struct {
	Count    int                  `json:"count"`
	Next     interface{}          `json:"next"`
	Previous interface{}          `json:"previous"`
	Results  []ResultsJobTemplate `json:"results"`
}
type RelatedJobTemplates struct {
	CreatedBy                    string `json:"created_by"`
	ModifiedBy                   string `json:"modified_by"`
	Labels                       string `json:"labels"`
	Inventory                    string `json:"inventory"`
	Project                      string `json:"project"`
	Organization                 string `json:"organization"`
	Credentials                  string `json:"credentials"`
	Jobs                         string `json:"jobs"`
	Schedules                    string `json:"schedules"`
	ActivityStream               string `json:"activity_stream"`
	Launch                       string `json:"launch"`
	WebhookKey                   string `json:"webhook_key"`
	WebhookReceiver              string `json:"webhook_receiver"`
	NotificationTemplatesStarted string `json:"notification_templates_started"`
	NotificationTemplatesSuccess string `json:"notification_templates_success"`
	NotificationTemplatesError   string `json:"notification_templates_error"`
	AccessList                   string `json:"access_list"`
	SurveySpec                   string `json:"survey_spec"`
	ObjectRoles                  string `json:"object_roles"`
	InstanceGroups               string `json:"instance_groups"`
	SliceWorkflowJobs            string `json:"slice_workflow_jobs"`
	Copy                         string `json:"copy"`
}
type Organization struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ScmType     string `json:"scm_type"`
}

type AdminRole struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
}

type UserCapabilities struct {
	Edit     bool `json:"edit"`
	Delete   bool `json:"delete"`
	Start    bool `json:"start"`
	Schedule bool `json:"schedule"`
	Copy     bool `json:"copy"`
}

type Credentials struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Cloud       bool   `json:"cloud"`
}

type Related struct {
	CreatedBy                    string `json:"created_by"`
	ModifiedBy                   string `json:"modified_by"`
	Labels                       string `json:"labels"`
	Inventory                    string `json:"inventory"`
	Project                      string `json:"project"`
	Organization                 string `json:"organization"`
	Credentials                  string `json:"credentials"`
	LastJob                      string `json:"last_job"`
	Jobs                         string `json:"jobs"`
	Schedules                    string `json:"schedules"`
	ActivityStream               string `json:"activity_stream"`
	Launch                       string `json:"launch"`
	WebhookKey                   string `json:"webhook_key"`
	WebhookReceiver              string `json:"webhook_receiver"`
	NotificationTemplatesStarted string `json:"notification_templates_started"`
	NotificationTemplatesSuccess string `json:"notification_templates_success"`
	NotificationTemplatesError   string `json:"notification_templates_error"`
	AccessList                   string `json:"access_list"`
	SurveySpec                   string `json:"survey_spec"`
	ObjectRoles                  string `json:"object_roles"`
	InstanceGroups               string `json:"instance_groups"`
	SliceWorkflowJobs            string `json:"slice_workflow_jobs"`
	Copy                         string `json:"copy"`
}
type SummaryFields struct {
	Organization     Organization     `json:"organization"`
	Inventory        Inventory        `json:"inventory"`
	Project          Project          `json:"project"`
	LastJob          LastJob          `json:"last_job"`
	LastUpdate       LastUpdate       `json:"last_update"`
	CreatedBy        CreatedBy        `json:"created_by"`
	ModifiedBy       ModifiedBy       `json:"modified_by"`
	ObjectRoles      ObjectRoles      `json:"object_roles"`
	UserCapabilities UserCapabilities `json:"user_capabilities"`
	Labels           Labels           `json:"labels"`
	RecentJobs       []RecentJobs     `json:"recent_jobs"`
	Credentials      []interface{}    `json:"credentials"`
}
type ResultsJobTemplate struct {
	ID                    int           `json:"id"`
	Type                  string        `json:"type"`
	URL                   string        `json:"url"`
	Related               Related       `json:"related,omitempty"`
	SummaryFields         SummaryFields `json:"summary_fields,omitempty"`
	Created               time.Time     `json:"created"`
	Modified              time.Time     `json:"modified"`
	Name                  string        `json:"name"`
	Description           string        `json:"description"`
	JobType               string        `json:"job_type"`
	Inventory             int           `json:"inventory"`
	Project               int           `json:"project"`
	Playbook              string        `json:"playbook"`
	ScmBranch             string        `json:"scm_branch"`
	Forks                 int           `json:"forks"`
	Limit                 string        `json:"limit"`
	Verbosity             int           `json:"verbosity"`
	ExtraVars             string        `json:"extra_vars"`
	JobTags               string        `json:"job_tags"`
	ForceHandlers         bool          `json:"force_handlers"`
	SkipTags              string        `json:"skip_tags"`
	StartAtTask           string        `json:"start_at_task"`
	Timeout               int           `json:"timeout"`
	UseFactCache          bool          `json:"use_fact_cache"`
	Organization          int           `json:"organization"`
	LastJobRun            interface{}   `json:"last_job_run"`
	LastJobFailed         bool          `json:"last_job_failed"`
	NextJobRun            interface{}   `json:"next_job_run"`
	Status                string        `json:"status"`
	HostConfigKey         string        `json:"host_config_key"`
	AskScmBranchOnLaunch  bool          `json:"ask_scm_branch_on_launch"`
	AskDiffModeOnLaunch   bool          `json:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch  bool          `json:"ask_variables_on_launch"`
	AskLimitOnLaunch      bool          `json:"ask_limit_on_launch"`
	AskTagsOnLaunch       bool          `json:"ask_tags_on_launch"`
	AskSkipTagsOnLaunch   bool          `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch    bool          `json:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch  bool          `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch  bool          `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch bool          `json:"ask_credential_on_launch"`
	SurveyEnabled         bool          `json:"survey_enabled"`
	BecomeEnabled         bool          `json:"become_enabled"`
	DiffMode              bool          `json:"diff_mode"`
	AllowSimultaneous     bool          `json:"allow_simultaneous"`
	CustomVirtualenv      interface{}   `json:"custom_virtualenv"`
	JobSliceCount         int           `json:"job_slice_count"`
	WebhookService        string        `json:"webhook_service"`
	WebhookCredential     interface{}   `json:"webhook_credential"`
}

type WorkflowJobTeplatesLaunchBody struct {
	ID        string `json:"id"`
	ExtraVars string `json:"extra_vars"`
}
