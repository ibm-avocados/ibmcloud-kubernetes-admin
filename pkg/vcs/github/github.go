package github

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/vcs"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

func CreateComment(comment vcs.GithubIssueComment) error {
	c, err := getCommentString(comment)
	if err != nil {
		log.Println("could not get comment string")
		return err
	}

	token := base64Encode(comment.GithubUser + ":" + comment.GithubToken)

	issueRepo := os.Getenv("GITHUB_ISSUE_REPO")

	base, owner, repo, err := processURL(issueRepo)
	if err != nil {
		return err
	}

	if err := notification.CreateComment(token, base, owner, repo, comment.IssueNumber, c); err != nil {
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
		return "", "", "", fmt.Errorf("malformed url: " + url)
	}
	return parts[0], parts[1], parts[2], nil
}
