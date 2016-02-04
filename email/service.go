package email

import (
	"net/http"
	"time"

	"github.com/RichardKnop/recall/config"
	"github.com/sendgrid/sendgrid-go"
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
	sg := sendgrid.NewSendGridClientWithApiKey(s.cnf.Sendgrid.APIKey)

	// Add *http.Client instance, so we can customise options such as the timeout
	sg.Client = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	// Construct the mail
	message := sendgrid.NewMail()
	message.SetSubject(email.Subject)
	for _, recipient := range email.Recipients {
		message.AddTo(recipient.Email)
		if recipient.Name != "" {
			message.AddToName(recipient.Name)
		}
	}
	message.SetFrom(email.From)
	message.SetText(email.Text)

	// And send the mail
	return sg.Send(message)
}
