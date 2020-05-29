package cron

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/notification"
)

func cloneGrantCluster(githubUser, githubToken, grantClusterURL string) error {
	if err := cmd("git",
		"clone", fmt.Sprintf("https://%s:%s@%s.git", githubUser, githubToken, grantClusterURL), "grant-cluster"); err != nil {
		return err
	}
	return nil
}

func login(apikey, org, space, resourceGroup, region string) error {
	if err := cmd("ibmcloud",
		"login", "--apikey", apikey, "-a", "https://cloud.ibm.com", "-r", region); err != nil {
		log.Println("login failed")
		return err
	}
	if err := cmd("ibmcloud",
		"target", "-o", org, "-s", space, "-g", resourceGroup); err != nil {
		log.Println("targetting org, space or region failed")
		return err
	}
	return nil
}

func deploy(apikey string, metadata *ibmcloud.AccountMetaData, schedule ibmcloud.Schedule) error {
	if err := login(apikey, metadata.Org, metadata.Space, schedule.ResourceGroupName, metadata.Region); err != nil {
		log.Println("could not login to ibmcloud")
		return err
	}

	if err := cloneGrantCluster(metadata.GithubUser, metadata.GithubToken, processURL(metadata.GrantClusterRepo)); err != nil {
		log.Println("could not clone grantcluster repo")
		return err
	}
	defer cleanUpFiles("grant-cluster")

	if err := cmd("./grant-cluster/scripts/deploy-app.sh"); err != nil {
		return err
	}

	comment, err := getCommentString(schedule, "templates/message.gotmpl")
	if err != nil {
		log.Println("could not get comment string")
		return err
	}

	token := base64.StdEncoding.EncodeToString([]byte(metadata.GithubUser + ":" + metadata.GithubToken))

	if err := notification.CreateComment(token, processURL(metadata.IssueRepo), schedule.GithubIssueNumber, comment); err != nil {
		return err
	}

	return nil
}

func getCommentString(schedule ibmcloud.Schedule, filename string) (string, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println("could not parse template files")
		return "", err
	}

	commentTemplate := template.Must(tmpl, err)
	buf := new(bytes.Buffer)

	if err := commentTemplate.Execute(buf, schedule); err != nil {
		log.Println("error executing comment template", err)
		return "", err
	}

	return buf.String(), nil
}

func cleanUp(apikey string, metadata *ibmcloud.AccountMetaData, schedule ibmcloud.Schedule) error {
	if err := login(apikey, metadata.Org, metadata.Space, schedule.ResourceGroupName, metadata.Region); err != nil {
		return err
	}

	if err := cloneGrantCluster(metadata.GithubUser, metadata.GithubToken, processURL(metadata.GrantClusterRepo)); err != nil {
		return err
	}
	defer cleanUpFiles("grant-cluster")
	if err := cmd("./grant-cluster/scripts/cleanup.sh"); err != nil {
		return err
	}
	return nil
}

// Cleanup folder
func cleanUpFiles(filepaths ...string) error {
	for _, filepath := range filepaths {
		fi, err := os.Stat(filepath)
		if err != nil {
			return err
		}
		mode := fi.Mode()
		if mode.IsDir() {
			if err := os.RemoveAll(filepath); err != nil {
				return err
			}
		} else if mode.IsRegular() {
			if err := os.Remove(filepath); err != nil {
				return err
			}
		}
	}
	return nil
}

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if output, err := cmd.Output(); err != nil {
		return err
	} else {
		fmt.Printf("%s\n", output)
	}
	return nil
}

func processURL(url string) string {
	if strings.Contains(url, "https://") {
		return strings.ReplaceAll(url, "https://", "")
	} else if strings.Contains(url, "http://") {
		return strings.ReplaceAll(url, "http://", "")
	}
	return url
}
