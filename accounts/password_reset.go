package accounts

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	// ErrPasswordResetNotFound ...
	ErrPasswordResetNotFound = errors.New("Password reset not found")
)

// FindPasswordResetByReference looks up a password reset by a reference
// only return the object if it's not expired
func (s *Service) FindPasswordResetByReference(reference string) (*PasswordReset, error) {
	// Fetch the password reset from the database
	passwordReset := new(PasswordReset)
	notFound := PasswordResetPreload(s.db).Where("reference = ?", reference).
		Where("expires_at > ?", time.Now().UTC()).First(passwordReset).RecordNotFound()

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
	notFound := PasswordResetPreload(s.db).Where("user_id = ?", user.ID).
		First(passwordReset).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrPasswordResetNotFound
	}

	return passwordReset, nil
}

func (s *Service) createPasswordReset(user *User) (*PasswordReset, error) {
	// Begin a transaction
	tx := s.db.Begin()

	// Soft delete old password resets
	if err := tx.Where("user_id = ?", user.ID).Delete(new(PasswordReset)).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new password reset
	passwordReset, err := NewPasswordReset(user, s.cnf.AppSpecific.PasswordResetLifetime)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Save the password reset to the database
	if err := tx.Create(passwordReset).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign related object
	passwordReset.User = user

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Send password reset email
	go func() {
		if err := s.sendPasswordResetEmail(passwordReset); err != nil {
			logger.Error(err)
		}
	}()

	return passwordReset, nil
}

func (s *Service) sendPasswordResetEmail(passwordReset *PasswordReset) error {
	passwordResetEmail, err := s.emailFactory.NewPasswordResetEmail(passwordReset)
	if err != nil {
		return fmt.Errorf("New password reset email error: %s", err)
	}

	// Try to send the password reset email
	if err := s.emailService.Send(passwordResetEmail); err != nil {
		return fmt.Errorf("Send email error: %s", err)
	}

	// If the email was sent successfully, update the email_sent flag
	now := gorm.NowFunc()
	if err := s.db.Model(passwordReset).UpdateColumns(PasswordReset{
		EmailTokenModel: EmailTokenModel{
			EmailSent:   true,
			EmailSentAt: &now,
			Model:       gorm.Model{UpdatedAt: now},
		},
	}).Error; err != nil {
		return nil
	}

	s.Notify()

	return nil
}
