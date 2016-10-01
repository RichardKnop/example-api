package accounts

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/RichardKnop/example-api/response"
	"github.com/RichardKnop/example-api/util"
)

// NewAccountAuthMiddleware creates a new AccountAuthMiddleware instance
func NewAccountAuthMiddleware(service ServiceInterface) *AccountAuthMiddleware {
	return &AccountAuthMiddleware{service: service}
}

// AccountAuthMiddleware takes the client ID and secret from the basic auth,
// authenticates the account and sets the account object on the request context
type AccountAuthMiddleware struct {
	service ServiceInterface
}

// ServeHTTP as per the negroni.Handler interface
func (m *AccountAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// HTTPS redirection
	err := util.NewSecure(m.service.GetConfig().IsDevelopment).Process(w, r)
	if err != nil {
		return
	}

	account, err := getClientCredentialsFromRequest(r, m.service)

	if err != nil || account == nil {
		// For security reasons, return a generic error message
		response.UnauthorizedError(w, ErrAccountAuthenticationRequired.Error())
		return
	}

	context.Set(r, AuthenticatedAccountKey, account)
	next(w, r)
}
