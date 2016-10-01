package accounts

import (
	"net/http"

	"github.com/gorilla/context"
)

type contextKey int

const (
	// AuthenticatedAccountKey ...
	AuthenticatedAccountKey contextKey = 0
	// AuthenticatedUserKey ...
	AuthenticatedUserKey contextKey = 1
)

// GetAuthenticatedAccount returns *Account from the request context
func GetAuthenticatedAccount(r *http.Request) (*Account, error) {
	val, ok := context.GetOk(r, AuthenticatedAccountKey)
	if !ok {
		return nil, ErrAccountAuthenticationRequired
	}

	authenticatedAccount, ok := val.(*Account)
	if !ok {
		return nil, ErrAccountAuthenticationRequired
	}

	return authenticatedAccount, nil
}

// GetAuthenticatedUser returns *User from the request context
func GetAuthenticatedUser(r *http.Request) (*User, error) {
	val, ok := context.GetOk(r, AuthenticatedUserKey)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	authenticatedUser, ok := val.(*User)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	return authenticatedUser, nil
}
