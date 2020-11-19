package github

import (
	"bytes"
	"html/template"
	"log"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func getCommentString(comment ibmcloud.GithubIssueComment, filename string) (string, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println("could not parse template files")
		return "", err
	}

	commentTemplate := template.Must(tmpl, err)
	buf := new(bytes.Buffer)

	if err := commentTemplate.Execute(buf, comment); err != nil {
		log.Println("error executing comment template", err)
		return "", err
	}

	return buf.String(), nil
}
