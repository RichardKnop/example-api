package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/models"
	"github.com/gorilla/context"
)

type contextKey string

const (
	// AuthenticatedClientKey ...
	AuthenticatedClientKey contextKey = "authenticated_client"
	// AuthenticatedUserKey ...
	AuthenticatedUserKey contextKey = "authenticated_user"
)

// GetAuthenticatedClient returns *OauthClient from the request context
func GetAuthenticatedClient(r *http.Request) (*models.OauthClient, error) {
	val, ok := context.GetOk(r, AuthenticatedClientKey)
	if !ok {
		return nil, ErrClientAuthenticationRequired
	}

	authenticatedClient, ok := val.(*models.OauthClient)
	if !ok {
		return nil, ErrClientAuthenticationRequired
	}

	return authenticatedClient, nil
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
