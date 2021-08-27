package github

import (
	"bytes"
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/vcs"
)

func getCommentString(comment vcs.GithubIssueComment) (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("can not get caller")
	}
	basePath := filepath.Dir(b)
	log.Println(basePath)
	file := filepath.Join(basePath, "templates/message.gotmpl")
	tmpl, err := template.ParseFiles(file)
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
