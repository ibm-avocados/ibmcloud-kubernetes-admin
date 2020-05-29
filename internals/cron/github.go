package cron

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

func createComment(schedule ibmcloud.Schedule, metadata *ibmcloud.AccountMetaData, templateFile string) error {
	comment, err := getCommentString(schedule, templateFile)
	if err != nil {
		log.Println("could not get comment string")
		return err
	}

	token := base64Encode(metadata.GithubUser + ":" + metadata.GithubToken)

	base, owner, repo, err := processURL(metadata.IssueRepo)
	if err != nil {
		return err
	}
	log.Println(base, owner, repo)

	if err := notification.CreateComment(token, base, owner, repo, schedule.GithubIssueNumber, comment); err != nil {
		log.Println("error posting comment to github", err)
		return err
	}
	return nil
}

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func processURL(url string) (string, string, string, error) {
	if strings.Contains(url, "https://") {
		url = strings.ReplaceAll(url, "https://", "")
	} else if strings.Contains(url, "http://") {
		url = strings.ReplaceAll(url, "http://", "")
	}
	parts := strings.Split(url, "/")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("malformed url")
	}
	return parts[0], parts[1], parts[2], nil
}
