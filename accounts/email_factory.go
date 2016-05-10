package accounts

import (
	"fmt"

	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/email"
)

var confirmationEmailTemplate = `
Hello %s,

Thank you for joining %s.

Please confirm your email: %s.

Kind Regards,

%s Team
`

var passwordResetEmailTemplate = `
Hello %s,

It seems you have forgotten your password.

You can set a new password here: %s.

Kind Regards,

%s Team
`

// EmailFactory facilitates construction of email.Email objects
type EmailFactory struct {
	cnf *config.Config
}

// NewEmailFactory starts a new emailFactory instance
func NewEmailFactory(cnf *config.Config) *EmailFactory {
	return &EmailFactory{cnf: cnf}
}

// NewConfirmationEmail returns a confirmation email
func (f *EmailFactory) NewConfirmationEmail(confirmation *Confirmation) *email.Email {
	// Define a greetings name for the user
	name := confirmation.User.GetName()
	if name == "" {
		name = "friend"
	}

	// App link
	appLink := fmt.Sprintf(
		"%s://%s",
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
	)

	// Confirmation link where the user can confirm his/her email
	link := fmt.Sprintf(
		"%s://%s/web/confirm-email/%s",
		f.cnf.Web.Scheme,
		f.cnf.Web.Host,
		confirmation.Reference,
	)

	// The email subject
	subject := fmt.Sprintf("Thank you for joining %s", f.cnf.Web.AppHost)

	// Replace placeholders in the email template
	emailText := fmt.Sprintf(
		confirmationEmailTemplate,
		name,
		appLink,
		link,
		appLink,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: confirmation.User.OauthUser.Username,
			Name:  confirmation.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.AppHost),
		Text: emailText,
	}
}

// NewPasswordResetEmail returns a password reset email
func (f *EmailFactory) NewPasswordResetEmail(passwordReset *PasswordReset) *email.Email {
	// Define a greetings name for the user
	name := passwordReset.User.GetName()
	if name == "" {
		name = "friend"
	}

	// App link
	appLink := fmt.Sprintf(
		"%s://%s",
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
	)

	// Password reset link where the user can set a new password
	link := fmt.Sprintf(
		"%s://%s/web/password-reset/%s",
		f.cnf.Web.Scheme,
		f.cnf.Web.Host,
		passwordReset.Reference,
	)

	// The email subject
	subject := fmt.Sprintf("Reset password for %s", f.cnf.Web.AppHost)

	// Replace placeholders in the email template
	emailText := fmt.Sprintf(
		passwordResetEmailTemplate,
		name,
		link,
		appLink,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: passwordReset.User.OauthUser.Username,
			Name:  passwordReset.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.AppHost),
		Text: emailText,
	}
}
