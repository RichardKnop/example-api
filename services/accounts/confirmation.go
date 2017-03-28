package accounts

import (
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/example-api/models"
	"github.com/jinzhu/gorm"
)

var (
	// ErrConfirmationNotFound ...
	ErrConfirmationNotFound = errors.New("Confirmation not found")
)

// FindConfirmationByReference looks up a confirmation by a reference
// only return the object if it's not expired
func (s *Service) FindConfirmationByReference(reference string) (*models.Confirmation, error) {
	// Fetch the invitation from the database
	confirmation := new(models.Confirmation)
	notFound := models.ConfirmationPreload(s.db).Where("reference = ?", reference).
		Where("expires_at > ?", time.Now().UTC()).First(confirmation).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrConfirmationNotFound
	}

	return confirmation, nil
}

// ConfirmUser confirms the user
func (s *Service) ConfirmUser(confirmation *models.Confirmation) error {
	// Begin a transaction
	tx := s.db.Begin()

	// Mark user as confirmed
	if err := tx.Model(confirmation.User).UpdateColumns(models.User{
		Confirmed: true,
		Model:     gorm.Model{UpdatedAt: time.Now().UTC()},
	}).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	// Soft delete the confirmation
	if err := tx.Delete(confirmation).Error; err != nil {
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

func (s *Service) sendConfirmationEmail(confirmation *models.Confirmation) error {
	confirmationEmail, err := s.emailFactory.NewConfirmationEmail(confirmation)
	if err != nil {
		return fmt.Errorf("New confirmation email error: %s", err)
	}

	// Try to send the confirmation email
	if err := s.emailService.Send(confirmationEmail); err != nil {
		return fmt.Errorf("Send email error: %s", err)
	}

	// If the email was sent successfully, update the email_sent flag
	now := gorm.NowFunc()
	if err := s.db.Model(confirmation).UpdateColumns(models.Confirmation{
		EmailTokenModel: models.EmailTokenModel{
			EmailSent:   true,
			EmailSentAt: &now,
			Model:       gorm.Model{UpdatedAt: now},
		},
	}).Error; err != nil {
		return err
	}

	s.Notify()

	return nil
}
