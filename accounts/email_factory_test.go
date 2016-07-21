package accounts_test

import (
	"io/ioutil"
	"testing"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
	"github.com/stretchr/testify/assert"
)

func TestNewConfirmationEmail(t *testing.T) {
	emailFactory := accounts.NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			AppScheme: "https",
			AppHost:   "example.com",
		},
		AppSpecific: config.AppSpecificConfig{
			CompanyName:  "Your Company Name",
			CompanyEmail: "contact@example.com",
		},
	})
	confirmation := &accounts.Confirmation{
		Reference: "some-reference",
		User: &accounts.User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
	}
	email, err := emailFactory.NewConfirmationEmail(confirmation)
	assert.NoError(t, err)
	assert.Equal(t, "Please confirm your email address", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "contact@example.com", email.From.Email)
	assert.Equal(t, "Your Company Name", email.From.Name)

	expectedPlain, err := ioutil.ReadFile("./accounts/test_templates/confirm_email.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedPlain), email.Text)

	expectedHTML, err := ioutil.ReadFile("./accounts/test_templates/confirm_email.html")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedHTML), email.HTML)
}

func TestNewPasswordResetEmail(t *testing.T) {
	emailFactory := accounts.NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			AppScheme: "https",
			AppHost:   "example.com",
		},
		AppSpecific: config.AppSpecificConfig{
			CompanyName:  "Your Company Name",
			CompanyEmail: "contact@example.com",
		},
	})
	passwordReset := &accounts.PasswordReset{
		Reference: "some-reference",
		User: &accounts.User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
	}
	email, err := emailFactory.NewPasswordResetEmail(passwordReset)
	assert.NoError(t, err)
	assert.Equal(t, "Reset password for example.com", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "contact@example.com", email.From.Email)
	assert.Equal(t, "Your Company Name", email.From.Name)

	expectedPlain, err := ioutil.ReadFile("./accounts/test_templates/password_reset_email.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedPlain), email.Text)

	expectedHTML, err := ioutil.ReadFile("./accounts/test_templates/password_reset_email.html")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedHTML), email.HTML)
}

func TestNewInvitationEmail(t *testing.T) {
	emailFactory := accounts.NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			AppScheme: "https",
			AppHost:   "example.com",
		},
		AppSpecific: config.AppSpecificConfig{
			CompanyName:  "Your Company Name",
			CompanyEmail: "contact@example.com",
		},
	})
	invitation := &accounts.Invitation{
		Reference: "some-reference",
		InvitedUser: &accounts.User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
		InvitedByUser: &accounts.User{
			OauthUser: &oauth.User{
				Username: "harold@finch",
			},
			FirstName: util.StringOrNull("Harold"),
			LastName:  util.StringOrNull("Finch"),
		},
	}
	email, err := emailFactory.NewInvitationEmail(invitation)
	assert.NoError(t, err)
	assert.Equal(t, "You have been invited to join example.com", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "harold@finch", email.From.Email)
	assert.Equal(t, "Harold Finch", email.From.Name)

	expectedPlain, err := ioutil.ReadFile("./accounts/test_templates/invitation_email.txt")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedPlain), email.Text)

	expectedHTML, err := ioutil.ReadFile("./accounts/test_templates/invitation_email.html")
	assert.NoError(t, err)
	assert.Equal(t, string(expectedHTML), email.HTML)
}
