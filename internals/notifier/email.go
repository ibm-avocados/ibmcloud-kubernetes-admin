package notifier

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

type Email struct {
	From        string
	Name        string
	Subject     string
	HtmlContent string
	Recipients  []string
}

func (e *Email) Send() error {
	m := mail.NewV3Mail()
	from := e.From
	name := e.Name
	sender := mail.NewEmail(name, from)
	m.SetFrom(sender)
	m.Subject = e.Subject
	p := mail.NewPersonalization()
	tos := make([]*mail.Email, len(e.Recipients))
	for i, recipient := range e.Recipients {
		tos[i] = mail.NewEmail(fmt.Sprintf("User %d", i+1), recipient)
	}

	p.AddTos(tos...)
	m.AddPersonalizations(p)

	html := mail.NewContent("text/html", e.HtmlContent)
	m.AddContent(html)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	res, err := client.Send(m)
	if err != nil {
		return err
	}
	log.Println("status code of email send", res.StatusCode)
	return nil
}

func NewEmail(content []byte) (*Email, error) {
	var email Email
	if err := json.Unmarshal(content, &email); err != nil {
		return nil, err
	}
	return &email, nil
}
