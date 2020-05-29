package cron

import (
	"fmt"
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func TestGetCommentString(t *testing.T) {
	request := ibmcloud.ClusterRequest{
		DataCenter:    "dal10",
		MachineType:   "4X16",
		MasterVersion: "iks16.7",
		WorkerNum:     1,
	}

	createRequest := ibmcloud.CreateClusterRequest{
		ClusterRequest: request,
	}

	schedule := ibmcloud.Schedule{
		Count:         "5",
		CreateRequest: createRequest,
		EventName:     "mofisapp",
		Password:      "ikslab",
	}

	comment, err := getCommentString(schedule, "../../templates/message.gotmpl")
	if err != nil {
		t.Errorf("there were error %v", err)
	}
	fmt.Println(comment)
}
