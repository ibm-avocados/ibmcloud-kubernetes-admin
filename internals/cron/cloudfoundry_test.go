package cron

import (
	"fmt"
	"testing"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func TestGetCommentString(t *testing.T) {
	request := ibmcloud.ScheduleRequest{
		DataCenters:   []string{"dal10"},
		MachineType:   "4X16",
		MasterVersion: "iks16.7",
		WorkerNum:     1,
	}

	scheduleRequest := ibmcloud.ScheduleClusterRequest{
		ScheduleRequest: request,
	}

	schedule := ibmcloud.Schedule{
		Count:           "5",
		ScheduleRequest: scheduleRequest,
		EventName:       "mofisapp",
		Password:        "ikslab",
	}

	comment, err := getCommentString(schedule, "../../templates/message.gotmpl")
	if err != nil {
		t.Errorf("there were error %v", err)
	}
	fmt.Println(comment)
}
