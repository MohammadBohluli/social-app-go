package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    sendgrid.NewSendClient(apiKey),
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template pasring
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return 0, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return 0, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return 0, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetires; i++ {
		resp, retryErr := m.client.Send(message)
		if retryErr != nil {
			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return resp.StatusCode, nil
	}

	return 0, fmt.Errorf("❌ failed to send email after %d attemps, error: %v", maxRetires, retryErr)

}
