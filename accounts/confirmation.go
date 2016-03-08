package accounts

import (
	"errors"
)

var (
	errConfirmationNotFound = errors.New("Confirmation not found")
)

// FindConfirmationByReference looks up a confirmation by a reference
func (s *Service) FindConfirmationByReference(reference string) (*Confirmation, error) {
	// Fetch the invitation from the database
	confirmation := new(Confirmation)
	notFound := s.db.Where("reference = ?", reference).
		Preload("User.OauthUser").First(confirmation).RecordNotFound()

	// Not found
	if notFound {
		return nil, errConfirmationNotFound
	}

	return confirmation, nil
}

// ConfirmUser sets User.Confirmed to true
func (s *Service) ConfirmUser(user *User) error {
	return s.db.Model(user).UpdateColumns(User{Confirmed: true}).Error
}
