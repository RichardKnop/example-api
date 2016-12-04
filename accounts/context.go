package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/models"
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
func GetAuthenticatedAccount(r *http.Request) (*models.Account, error) {
	val, ok := context.GetOk(r, AuthenticatedAccountKey)
	if !ok {
		return nil, ErrAccountAuthenticationRequired
	}

	authenticatedAccount, ok := val.(*models.Account)
	if !ok {
		return nil, ErrAccountAuthenticationRequired
	}

	return authenticatedAccount, nil
}

// GetAuthenticatedUser returns *User from the request context
func GetAuthenticatedUser(r *http.Request) (*models.User, error) {
	val, ok := context.GetOk(r, AuthenticatedUserKey)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	authenticatedUser, ok := val.(*models.User)
	if !ok {
		return nil, ErrUserAuthenticationRequired
	}

	return authenticatedUser, nil
}
