package accounts

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/RichardKnop/example-api/util/response"
	"github.com/RichardKnop/example-api/util"
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
	// HTTPS redirection
	err := util.NewSecure(m.service.GetConfig().IsDevelopment).Process(w, r)
	if err != nil {
		return
	}

	account, user, err := getMixedCredentialsFromRequest(r, m.service)

	if err != nil {
		// For security reasons, return a generic error message
		response.UnauthorizedError(w, ErrAccountOrUserAuthenticationRequired.Error())
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
