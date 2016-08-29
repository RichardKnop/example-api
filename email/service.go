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
func (s *Service) Send(m *Message) error {
	// Construct the mail
	message := new(mail.SGMailV3)
	for _, recipient := range m.Recipients {
		message.SetFrom(&mail.Email{Address: recipient.Address, Name: recipient.Name})
	}
	message.Subject = m.Subject
	p := mail.NewPersonalization()
	for _, recipient := range m.Recipients {
		p.AddTos(&mail.Email{Address: recipient.Address, Name: recipient.Name})
	}
	message.AddPersonalizations(p)
	content := mail.NewContent("text/plain", m.Text)
	message.AddContent(content)

	// And send the mail
	request := sendgrid.GetRequest(
		s.cnf.Sendgrid.APIKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(message)
	_, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	return nil
}
