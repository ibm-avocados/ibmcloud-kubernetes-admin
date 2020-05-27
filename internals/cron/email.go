package cron

import (
	"bytes"
	"html/template"
	"log"
)

func getEmailBody(data EmailData) (string, error) {
	tmpl, err := template.ParseFiles("templates/email.gohtml")
	if err != nil {
		log.Println("could not parse file", err)
		return "", err
	}
	htmlTemplate := template.Must(tmpl, err)
	buf := new(bytes.Buffer)

	if err := htmlTemplate.Execute(buf, data); err != nil {
		log.Println("could not parse file", err)
		return "", err
	}

	return buf.String(), nil
}
