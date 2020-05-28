package cron

import (
	"os"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func setEnvs(accountID string, metadata *ibmcloud.AccountMetaData, schedule ibmcloud.Schedule) error {

	if err := os.Setenv("EVENT_NAME", schedule.EventName); err != nil {
		return err
	}
	if err := os.Setenv("PASSWORD", schedule.Password); err != nil {
		return err
	}
	if err := os.Setenv("RESOURCE_GROUP_NAME", schedule.ResourceGroupName); err != nil {
		return err
	}
	if err := os.Setenv("APP_HOSTNAME", schedule.EventName); err != nil {
		return err
	}
	if err := os.Setenv("ACCOUNT", accountID); err != nil {
		return err
	}
	if err := os.Setenv("FILTER_TAG", schedule.EventName); err != nil {
		return err
	}
	if err := os.Setenv("ACCESS_GROUP_NAME", metadata.AccessGroup); err != nil {
		return err
	}
	return nil
}
