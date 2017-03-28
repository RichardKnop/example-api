package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/util/response"
	"github.com/gorilla/context"
)

var (
	// ErrClientOrUserAuthenticationRequired ...
	ErrClientOrUserAuthenticationRequired = errors.New("Client or user authentication required")
)

// NewMixedAuthMiddleware creates a new MixedAuthMiddleware instance
func NewMixedAuthMiddleware(service ServiceInterface) *MixedAuthMiddleware {
	return &MixedAuthMiddleware{service: service}
}

// MixedAuthMiddleware looks for either client or user auth and accepts both
type MixedAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *MixedAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	client, user, err := getMixedCredentialsFromRequest(r, m.service)

	if err != nil {
		// For security reasons, return a generic error message
		response.UnauthorizedError(w, ErrClientOrUserAuthenticationRequired.Error())
		return
	}

	if client != nil {
		context.Set(r, AuthenticatedClientKey, client)
	}

	if user != nil {
		context.Set(r, AuthenticatedUserKey, user)
	}

	next(w, r)
}
