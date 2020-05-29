package cron

import (
	"encoding/base64"
	"os"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

func createComment(issue, comment string) error {
	githubUser := os.Getenv("GITHUB_USER")
	githubToken := os.Getenv("GITHUB_TOKEN")
	token := "Basic " + base64Encode(githubUser+":"+githubToken)
	if err := notification.CreateComment(token, "", issue, comment); err != nil {
		return err
	}
	return nil
}

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}