package accounts

import (
	"net/http"

	"github.com/RichardKnop/recall/response"
	"github.com/RichardKnop/recall/util"
	"github.com/gorilla/context"
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
	// HTTPS redirection
	err := util.NewSecure(m.service.GetConfig().IsDevelopment).Process(w, r)
	if err != nil {
		return
	}

	account, user, err := getCredentialsFromRequest(r, m.service)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	if account != nil {
		context.Set(r, AuthenticatedAccountKey, account)
	}

	if user != nil {
		context.Set(r, AuthenticatedUserKey, user)
	} else {
		response.UnauthorizedError(w, ErrUserAuthenticationRequired.Error())
		return
	}

	next(w, r)
}

// NewOptionalUserAuthMiddleware creates a new OptionalUserAuthMiddleware instance
func NewOptionalUserAuthMiddleware(service ServiceInterface) *OptionalUserAuthMiddleware {
	return &OptionalUserAuthMiddleware{service: service}
}

// OptionalUserAuthMiddleware takes the bearer token from the Authorization header,
// authenticates the user and sets the user object on the request context. If it
// cannot find it, it just continues
type OptionalUserAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *OptionalUserAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// HTTPS redirection
	err := util.NewSecure(m.service.GetConfig().IsDevelopment).Process(w, r)
	if err != nil {
		return
	}

	account, user, err := getCredentialsFromRequest(r, m.service)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	if account != nil {
		context.Set(r, AuthenticatedAccountKey, account)
	}

	if user != nil {
		context.Set(r, AuthenticatedUserKey, user)
	}

	next(w, r)
}
