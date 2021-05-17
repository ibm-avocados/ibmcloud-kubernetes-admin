package github

import (
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func TestCommentString(t *testing.T) {
	issue := ibmcloud.GithubIssueComment{
		IssueNumber: "",
		EventName:   "ibmcloudevent",
		Password:    "password",
		AccountID:   "1234",
		GithubUser:  "user",
		GithubToken: "",
		ClusterRequest: ibmcloud.GithubIssueClusterRequest{
			Count:      10,
			Type:       "kubernetes",
			ErrorCount: 0,
			Regions:    "dal-10,dal-12",
		},
	}
	_, err := getCommentString(issue, "../../templates/message.gotmpl")
	if err != nil {
		t.Fatal(err)
	}
}
