package email

import (
	"github.com/RichardKnop/example-api/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Service struct keeps config object to avoid passing it around
type Service struct {
	cnf *config.Config
}

// NewService starts a new Service instance
func NewService(cnf *config.Config) *Service {
	return &Service{cnf: cnf}
}

// Send sends email using sendgrid
func (s *Service) Send(email *Email) error {
	// Construct the mail
	m := new(mail.SGMailV3)
	m.SetFrom(&mail.Email{Address: email.From.Email, Name: email.From.Name})
	m.Subject = email.Subject
	p := mail.NewPersonalization()
	for _, recipient := range email.Recipients {
		p.AddTos(&mail.Email{Address: recipient.Email, Name: recipient.Name})
	}
	m.AddPersonalizations(p)
	content := mail.NewContent("text/plain", email.Text)
	m.AddContent(content)

	// And send the mail
	request := sendgrid.GetRequest(
		s.cnf.Sendgrid.APIKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	return nil
}
