package github

import (
	"os"
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vcs"
)

func TestCreateComment(t *testing.T) {
	issue := vcs.GithubIssueComment{
		IssueNumber: os.Getenv("TEST_GITHUB_REPO"),
		EventName:   "functiontest",
		Password:    "password",
		AccountID:   "1234",
		GithubUser:  "Mofizur-Rahman",
		GithubToken: os.Getenv("TEST_GITHUB_TOKEN"),
		ClusterRequest: vcs.GithubIssueClusterRequest{
			Count:      10,
			Type:       "kubernetes",
			ErrorCount: 0,
			Regions:    "dal-10,dal-12",
		},
	}
	err := CreateComment(issue)
	if err != nil {
		t.Fatal(err)
	}
}
