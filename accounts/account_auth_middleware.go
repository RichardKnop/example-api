package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/response"
	"github.com/RichardKnop/example-api/util"
	"github.com/gorilla/context"
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

	account, _, err := getCredentialsFromRequest(r, m.service)

	if err != nil {
		response.UnauthorizedError(w, ErrAccountAuthenticationRequired.Error())
		return
	}

	if account != nil {
		context.Set(r, AuthenticatedAccountKey, account)
		next(w, r)
	} else {
		response.UnauthorizedError(w, ErrAccountAuthenticationRequired.Error())
		return
	}
}
