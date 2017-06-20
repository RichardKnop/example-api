package accounts

import (
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(o *models.Confirmation) (*email.Message, error)
	NewInvitationEmail(o *models.Invitation) (*email.Message, error)
	NewPasswordResetEmail(o *models.PasswordReset) (*email.Message, error)
}
