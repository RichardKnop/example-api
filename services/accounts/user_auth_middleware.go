package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/util/response"
	"github.com/gorilla/context"
)

var (
	// ErrUserAuthenticationRequired ...
	ErrUserAuthenticationRequired = errors.New("User authentication required")
)

// NewUserAuthMiddleware creates a new UserAuthMiddleware instance
func NewUserAuthMiddleware(service ServiceInterface) *UserAuthMiddleware {
	return &UserAuthMiddleware{service: service}
}

// UserAuthMiddleware takes the bearer token from the Authorization header,
// authenticates the user and sets the user object on the request context. If it
// cannot find it then it throws an unauthorized error
type UserAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *UserAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, err := getUserCredentialsFromRequest(r, m.service)
	if err != nil || user == nil {
		// For security reasons, return a generic error message
		response.UnauthorizedError(w, ErrUserAuthenticationRequired.Error())
		return
	}

	context.Set(r, AuthenticatedUserKey, user)
	next(w, r)
}

// NewOptionalAuthMiddleware creates a new OptionalAuthMiddleware instance
func NewOptionalAuthMiddleware(service ServiceInterface) *OptionalAuthMiddleware {
	return &OptionalAuthMiddleware{service: service}
}

// OptionalAuthMiddleware takes the bearer token from the Authorization header,
// authenticates the client and/or user and sets the client and/or user objects
// on the request context. If it cannot find it, it just continues
type OptionalAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *OptionalAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Optional user auth
	user, err := getUserCredentialsFromRequest(r, m.service)
	if err == nil && user != nil {
		context.Set(r, AuthenticatedUserKey, user)
		context.Set(r, AuthenticatedClientKey, user.OauthClient)
	}

	// Optional client auth
	client, err := getClientCredentialsFromRequest(r, m.service)
	if err == nil && client != nil {
		context.Set(r, AuthenticatedClientKey, client)
	}

	next(w, r)
}
