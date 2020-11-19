package cron

import (
	"fmt"
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func TestGetCommentString(t *testing.T) {
	request := ibmcloud.GithubIssueClusterRequest{
		Count:      10,
		ErrorCount: 0,
		Regions:    "dal10,dal12",
		Type:       "Openshift",
	}

	comment := ibmcloud.GithubIssueComment{
		EventName:      "mofisapp",
		Password:       "ikslab",
		ClusterRequest: request,
		AccountID:      "64cf1b006290",
		IssueNumber:    "1234",
	}

	c, err := getCommentString(comment, "../../templates/message.gotmpl")
	if err != nil {
		t.Errorf("there were error %v", err)
	}
	fmt.Println(c)
}
