package notification

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Email(subject, htmlContent string, recipients ...string) error {
	m := mail.NewV3Mail()
	from := os.Getenv("ADMIN_FROM_EMAIL")
	name := "IBMCloud Kube Admin"
	e := mail.NewEmail(name, from)
	m.SetFrom(e)
	m.Subject = subject
	p := mail.NewPersonalization()
	tos := make([]*mail.Email, len(recipients))
	for i, recipient := range recipients {
		tos[i] = mail.NewEmail(fmt.Sprintf("User %d", i+1), recipient)
	}
	p.AddTos(tos...)

	m.AddPersonalizations(p)

	html := mail.NewContent("text/html", htmlContent)
	m.AddContent(html)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	res, err := client.Send(m)
	if err != nil {
		return err
	}
	log.Println("status code of email", res.StatusCode)
	return nil
}

func EmailAdmin(subject, htmlContent string) error {
	adminEmail := os.Getenv("ADMIN_TO_EMAIL")
	if adminEmail == "" {
		return fmt.Errorf("no admin email provided")
	}
	return Email(subject, htmlContent, strings.Split(adminEmail, ",")...)
}
