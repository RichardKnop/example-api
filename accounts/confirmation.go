package accounts

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	// ErrConfirmationNotFound ...
	ErrConfirmationNotFound = errors.New("Confirmation not found")
)

// FindConfirmationByReference looks up a confirmation by a reference
func (s *Service) FindConfirmationByReference(reference string) (*Confirmation, error) {
	// Fetch the invitation from the database
	confirmation := new(Confirmation)
	notFound := ConfirmationPreload(s.db).Where("reference = ?", reference).
		First(confirmation).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrConfirmationNotFound
	}

	return confirmation, nil
}

// ConfirmUser sets confirmed flag to true
func (s *Service) ConfirmUser(user *User) error {
	return s.db.Model(user).UpdateColumns(User{
		Confirmed: true,
		Model:     gorm.Model{UpdatedAt: time.Now()},
	}).Error
}
