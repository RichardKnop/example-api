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
			Host:      "api.recall",
			AppScheme: "https",
			AppHost:   "recall",
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

	assert.Equal(t, "Thank you for joining recall", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@recall", email.From.Email)
	assert.Equal(t, "NOREPLY recall", email.From.Name)

	expectedText := `
Hello John Reese,

Thank you for joining https://recall.

Please confirm your email: https://api.recall/web/confirm-email/some-reference.

Kind Regards,

https://recall Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewInvitationEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.recall",
			AppScheme: "https",
			AppHost:   "recall",
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

	assert.Equal(t, "You have been invited to join recall", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "harold@finch", email.From.Email)
	assert.Equal(t, "Harold Finch", email.From.Name)

	expectedText := `
Hello John Reese,

You have been invited to join https://recall by Harold Finch.

Follow this link to set your password please: https://api.recall/web/confirm-invitation/some-reference.

Kind Regards,

https://recall Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewPasswordResetEmail(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.recall",
			AppScheme: "https",
			AppHost:   "recall",
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

	assert.Equal(t, "Reset password for recall", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@recall", email.From.Email)
	assert.Equal(t, "NOREPLY recall", email.From.Name)

	expectedText := `
Hello John Reese,

It seems you have forgotten your password.

You can set a new password here: https://api.recall/web/password-reset/some-reference.

Kind Regards,

https://recall Team
`
	assert.Equal(t, expectedText, email.Text)
}
