package email

import (
	"fmt"

	"github.com/RichardKnop/example-api/config"
	"github.com/aymerick/douceur/inliner"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// Service struct keeps config object to avoid passing it around
type Service struct {
	cnf *config.Config
}

// NewService starts a new Service instance
func NewService(cnf *config.Config) *Service {
	return &Service{cnf: cnf}
}

// Send sends email message using mailgun
func (s *Service) Send(m *Message) error {
	// Format recipients
	var formattedRecipients []string
	for _, recipient := range m.Recipients {
		formattedAddress, err := recipient.Format()
		if err != nil {
			return err
		}
		formattedRecipients = append(formattedRecipients, formattedAddress)
	}

	// Format the sender
	formattedSender, err := m.From.Format()
	if err != nil {
		return err
	}

	// Create email message
	message := mailgun.NewMessage(
		formattedSender,
		m.Subject,
		m.Text,
		formattedRecipients...,
	)

	// Optionally set HTML body
	if m.HTML != "" {
		htmlWithInlineCSS, err := inliner.Inline(m.HTML)
		if err != nil {
			return fmt.Errorf("CSS inliner error: %s", err.Error())
		}
		message.SetHtml(htmlWithInlineCSS)
	}

	// TODO - do we need to return other values than error here?
	mg := mailgun.NewMailgun(
		s.cnf.Mailgun.Domain,
		s.cnf.Mailgun.APIKey,
		s.cnf.Mailgun.PublicAPIKey,
	)
	_, _, err = mg.Send(message)
	return err
}
