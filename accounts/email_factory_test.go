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
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
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

	assert.Equal(t, "Thank you for joining pingli.st", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Thank you for joining https://pingli.st.

Please confirm your email: https://api.pingli.st/web/confirm-email/some-reference.

Kind Regards,

https://pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewInvitationEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})
	invitation := &Invitation{
		Reference: "some-reference",
		InvitedUser: &User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
		InvitedByUser: &User{
			OauthUser: &oauth.User{
				Username: "harold@finch",
			},
			FirstName: util.StringOrNull("Harold"),
			LastName:  util.StringOrNull("Finch"),
		},
	}
	email := emailFactory.NewInvitationEmail(invitation)

	assert.Equal(t, "You have been invited to join pingli.st", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "harold@finch", email.From)

	expectedText := `
Hello John Reese,

You have been invited to join https://pingli.st by Harold Finch.

Follow this link to set your password please: https://api.pingli.st/web/confirm-invitation/some-reference.

Kind Regards,

https://pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewPasswordResetEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
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

	assert.Equal(t, "Reset password for pingli.st", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

It seems you have forgotten your password.

You can set a new password here: https://api.pingli.st/web/password-reset/some-reference.

Kind Regards,

https://pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}
