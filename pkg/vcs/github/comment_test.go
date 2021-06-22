package github

import (
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vcs"
)

func TestCommentString(t *testing.T) {
	issue := vcs.GithubIssueComment{
		IssueNumber: "",
		EventName:   "ibmcloudevent",
		Password:    "password",
		AccountID:   "1234",
		GithubUser:  "user",
		GithubToken: "",
		ClusterRequest: vcs.GithubIssueClusterRequest{
			Count:      10,
			Type:       "kubernetes",
			ErrorCount: 0,
			Regions:    "dal-10,dal-12",
		},
	}
	_, err := getCommentString(issue)
	if err != nil {
		t.Fatal(err)
	}
}
