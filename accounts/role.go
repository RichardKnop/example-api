package accounts

import (
	"errors"
)

var (
	// ErrRoleNotFound ...
	ErrRoleNotFound = errors.New("Role not found")
)

// findRoleByID looks up a role by ID and returns it
func (s *Service) findRoleByID(id string) (*Role, error) {
	role := new(Role)
	if s.db.Where("id = ?", id).First(role).RecordNotFound() {
		return nil, ErrRoleNotFound
	}
	return role, nil
}
