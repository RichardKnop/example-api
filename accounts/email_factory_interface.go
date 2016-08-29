package accounts

import (
	"github.com/RichardKnop/example-api/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(confirmation *Confirmation) (*email.Message, error)
	NewInvitationEmail(invitation *Invitation) (*email.Message, error)
	NewPasswordResetEmail(passwordReset *PasswordReset) (*email.Message, error)
}
