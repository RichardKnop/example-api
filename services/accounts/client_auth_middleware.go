package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/util/response"
	"github.com/gorilla/context"
)

var (
	// ErrClientAuthenticationRequired ...
	ErrClientAuthenticationRequired = errors.New("Client authentication required")
)

// NewClientAuthMiddleware creates a new ClientAuthMiddleware instance
func NewClientAuthMiddleware(service ServiceInterface) *ClientAuthMiddleware {
	return &ClientAuthMiddleware{service: service}
}

// ClientAuthMiddleware takes the client ID and secret from the basic auth,
// authenticates the client and sets the client object on the request context
type ClientAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *ClientAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	client, err := getClientCredentialsFromRequest(r, m.service)

	if err != nil || client == nil {
		// For security reasons, return a generic error message
		response.UnauthorizedError(w, ErrClientAuthenticationRequired.Error())
		return
	}

	context.Set(r, AuthenticatedClientKey, client)
	next(w, r)
}
