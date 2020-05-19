package notification

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Email(subject, textContent, htmlContent string, recipients []string) error {
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

	text := mail.NewContent("text/plain", textContent)
	m.AddContent(text)

	html := mail.NewContent("text/html", htmlContent)
	m.AddContent(html)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	res, err := client.Send(m)
	if err != nil {
		return err
	}
	log.Println(res.StatusCode)
	return nil
}
