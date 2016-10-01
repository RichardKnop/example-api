package accounts

import (
	"github.com/RichardKnop/example-api/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(o *Confirmation) (*email.Message, error)
	NewInvitationEmail(o *Invitation) (*email.Message, error)
	NewPasswordResetEmail(o *PasswordReset) (*email.Message, error)
}
