package accounts

import (
	"testing"

	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/oauth"
	"github.com/RichardKnop/recall/util"
	"github.com/stretchr/testify/assert"
)

func TestNewConfirmationEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.example.com",
			AppScheme: "https",
			AppHost:   "example.com",
		},
	})
	confirmation := &Confirmation{
		Reference: "some-reference",
		User: &User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
	}
	email := emailFactory.NewConfirmationEmail(confirmation)

	assert.Equal(t, "Thank you for joining example.com", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@example.com", email.From)

	expectedText := `
Hello John Reese,

Thank you for joining https://example.com.

Please confirm your email: https://api.example.com/web/confirm-email/some-reference.

Kind Regards,

https://example.com Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewPasswordResetEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.example.com",
			AppScheme: "https",
			AppHost:   "example.com",
		},
	})
	passwordReset := &PasswordReset{
		Reference: "some-reference",
		User: &User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
	}
	email := emailFactory.NewPasswordResetEmail(passwordReset)

	assert.Equal(t, "Reset password for example.com", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@example.com", email.From)

	expectedText := `
Hello John Reese,

It seems you have forgotten your password.

You can set a new password here: https://api.example.com/web/password-reset/some-reference.

Kind Regards,

https://example.com Team
`
	assert.Equal(t, expectedText, email.Text)
}
