package accounts

import (
	"net/http"

	"github.com/RichardKnop/recall/response"
	"github.com/RichardKnop/recall/util"
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

	// First, try to get the bearer token
	token, err := util.ParseBearerToken(r)
	if err == nil {
		// Authenticate
		oauthAccessToken, err := m.service.GetOauthService().Authenticate(string(token))
		if err != nil {
			// For security reasons, return a general error message
			response.UnauthorizedError(w, errAccountAuthenticationRequired.Error())
			return
		}

		// Fetch the account
		account, err := m.service.FindAccountByOauthClientID(oauthAccessToken.Client.ID)
		if err != nil {
			// For security reasons, return a general error message
			response.UnauthorizedError(w, errAccountAuthenticationRequired.Error())
			return
		}

		context.Set(r, authenticatedAccountKey, account)
		next(w, r)
		return
	}

	// Second, try to get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errAccountAuthenticationRequired.Error())
		return
	}

	// Authenticate the client
	oauthClient, err := m.service.GetOauthService().AuthClient(clientID, secret)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errAccountAuthenticationRequired.Error())
		return
	}

	// Fetch the account
	account, err := m.service.FindAccountByOauthClientID(oauthClient.ID)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errAccountAuthenticationRequired.Error())
		return
	}

	context.Set(r, authenticatedAccountKey, account)
	next(w, r)
}
