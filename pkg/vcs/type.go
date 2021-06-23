package vcs

type GithubIssueComment struct {
	IssueNumber    string                    `json:"issueNumber"`
	EventName      string                    `json:"eventName"`
	Password       string                    `json:"password"`
	AccountID      string                    `json:"accountID"`
	GithubUser     string                    `json:"githubUser"`
	GithubToken    string                    `json:"githubToken"`
	ClusterRequest GithubIssueClusterRequest `json:"clusterRequest"`
}

type GithubIssueClusterRequest struct {
	Count      int    `json:"count"`
	Type       string `json:"type"`
	ErrorCount int    `json:"errorCount"`
	Regions    string `json:"regions"`
}
