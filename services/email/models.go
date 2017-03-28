package email

import (
	"errors"
	"fmt"
)

// Email represents an email address with a name
type Email struct {
	Address string
	Name    string
}

// Format returns a properly formatted email address
func (e *Email) Format() (string, error) {
	if e.Address == "" {
		return "", errors.New("Email address cannot be empty")
	}
	if e.Name == "" {
		return e.Address, nil
	}
	return fmt.Sprintf("%s <%s>", e.Name, e.Address), nil
}

// Recipient represents a single email recipient
type Recipient struct {
	Email
}

// Sender represents an email sender
type Sender struct {
	Email
}

// Message represents an email message
type Message struct {
	Subject    string
	Recipients []*Recipient
	From       *Sender
	Text       string
	HTML       string
}
