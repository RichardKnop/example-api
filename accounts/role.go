package accounts

import (
	"errors"
)

var (
	errRoleNotFound = errors.New("Role not found")
)

// findRoleByName looks up a role by name and returns it
func (s *Service) findRoleByName(name string) (*Role, error) {
	role := new(Role)
	if s.db.Where(Role{
		Name: name,
	}).First(role).RecordNotFound() {
		return nil, errRoleNotFound
	}
	return role, nil
}
