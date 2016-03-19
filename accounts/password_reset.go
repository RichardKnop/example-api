package accounts

import (
	"errors"

	"github.com/RichardKnop/recall/util"
)

var (
	// ErrPasswordResetNotFound ...
	ErrPasswordResetNotFound = errors.New("Password reset not found")
)

// FindPasswordResetByReference looks up a password reset by a reference
func (s *Service) FindPasswordResetByReference(reference string) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	notFound := s.db.Where("reference = ?", reference).
		Preload("User.OauthUser").First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}

// ResetPassword sets a new password and deletes the password reset record
func (s *Service) ResetPassword(passwordReset *PasswordReset, password string) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Set the new password
	err := s.oauthService.SetPasswordTx(tx, passwordReset.User.OauthUser, password)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Soft delete the password reset
	if err := tx.Delete(passwordReset).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}

// findUserPasswordReset returns the first password reset for a user
func (s *Service) findUserPasswordReset(user *User) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	notFound := s.db.Where(PasswordReset{
		UserID: util.PositiveIntOrNull(int64(user.ID)),
	}).Preload("User.OauthUser").First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}

func (s *Service) createPasswordReset(user *User) (*PasswordReset, error) {
	var passwordReset *PasswordReset

	// Does the user have an open password reset?
	passwordReset, err := s.findUserPasswordReset(user)
	if err != nil {
		// Create a new password reset
		passwordReset = newPasswordReset(user)
	}

	// Save the password reset to the database
	if err := s.db.Create(passwordReset).Error; err != nil {
		return nil, err
	}

	return passwordReset, nil
}
