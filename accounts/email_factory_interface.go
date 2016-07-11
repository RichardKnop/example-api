package accounts

import (
	"github.com/RichardKnop/recall/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(confirmation *Confirmation) (*email.Email, error)
	NewInvitationEmail(invitation *Invitation) (*email.Email, error)
	NewPasswordResetEmail(passwordReset *PasswordReset) (*email.Email, error)
}
