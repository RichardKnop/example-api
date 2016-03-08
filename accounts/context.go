package accounts

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
)

type contextKey int

const (
	authenticatedAccountKey contextKey = 0
	authenticatedUserKey    contextKey = 1
)

var (
	// ErrAccountAuthenticationRequired ...
	ErrAccountAuthenticationRequired = errors.New("Account authentication required")
	// ErrUserAuthenticationRequired ...
	ErrUserAuthenticationRequired = errors.New("User authentication required")
)

// GetAuthenticatedAccount returns *Account from the request context
func GetAuthenticatedAccount(r *http.Request) (*Account, error) {
	val, ok := context.GetOk(r, authenticatedAccountKey)
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
	val, ok := context.GetOk(r, authenticatedUserKey)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	authenticatedUser, ok := val.(*User)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	return authenticatedUser, nil
}
