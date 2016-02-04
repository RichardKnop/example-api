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
// authenticates the user and sets the user object on the request context
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

	// Get the bearer token
	token, err := util.ParseBearerToken(r)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errUserAuthenticationRequired.Error())
		return
	}

	// Authenticate
	oauthAccessToken, err := m.service.GetOauthService().Authenticate(string(token))
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errUserAuthenticationRequired.Error())
		return
	}

	// Access token has no user, this probably means client credentials grant
	if oauthAccessToken.User == nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errUserAuthenticationRequired.Error())
		return
	}

	// Fetch the user account
	user, err := m.service.FindUserByOauthUserID(oauthAccessToken.User.ID)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errUserAuthenticationRequired.Error())
		return
	}

	context.Set(r, authenticatedUserKey, user)

	next(w, r)
}
