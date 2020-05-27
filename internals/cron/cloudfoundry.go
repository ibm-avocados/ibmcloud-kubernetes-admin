package cron

import (
	"fmt"
	"os"
	"os/exec"
)

func deploy(apikey, org, space, resourceGroup, region string) error {
	githubUser := os.Getenv("GITHUB_USER")
	githubToken := os.Getenv("GITHUB_TOKEN")
	grantClusterURL := os.Getenv("GRANT_CLUSTER_REPO_URL")

	if err := cmd("ibmcloud",
		"login", "--apikey", apikey, "-a", "https://cloud.ibm.com", "-r", region); err != nil {
		return err
	}
	if err := cmd("ibmcloud",
		"target", "-o", org, "-s", space, "-g", resourceGroup); err != nil {
		return err
	}
	if err := cmd("git",
		"clone", fmt.Sprintf("https://%s:%s@%s.git", githubUser, githubToken, grantClusterURL)); err != nil {
		return err
	}

	defer cleanup("grant-cluster")

	if err := cmd("./grant-cluster/scripts/deploy-app.sh"); err != nil {
		return err
	}

	return nil
}

// Cleanup folder
func cleanup(filepaths ...string) error {
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
