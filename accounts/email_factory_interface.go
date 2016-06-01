package accounts

import (
	"github.com/RichardKnop/recall/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewConfirmationEmail(confirmation *Confirmation) *email.Email
	NewInvitationEmail(invitation *Invitation) *email.Email
	NewPasswordResetEmail(passwordReset *PasswordReset) *email.Email
}
