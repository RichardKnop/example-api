package facebook

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/oauth"
	"github.com/RichardKnop/recall/response"
)

var (
	// ErrAccountMismatch ...
	ErrAccountMismatch = errors.New("Account mismatch")
)

// Handles requests to login with Facebook access token (POST /v1/facebook/login)
func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated account from the request context
	authenticatedAccount, err := accounts.GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the scope string
	scope, err := s.GetAccountsService().GetOauthService().GetScope(r.Form.Get("scope"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user data from facebook
	resp, err := s.adapter.GetMe(r.Form.Get("access_token"))
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Initialise variables with from facebook
	var (
		facebookID = fmt.Sprintf("%s", resp["id"])
		email      = fmt.Sprintf("%s", resp["email"])
		firstName  = fmt.Sprintf("%s", resp["first_name"])
		lastName   = fmt.Sprintf("%s", resp["last_name"])
		user       *accounts.User
	)

	logger.Info("Fetched Facebook user's data")
	logger.Infof("%v", resp)

	// There is an edge case where Facebook does not return a valid email
	// User could have registered with a phone number or have an unconfirmed
	// email address. In such rare case, default to {facebook_id}@facebook.com
	if resp["email"] == nil || email == "%!s(<nil>)" {
		email = fmt.Sprintf("%s@facebook.com", facebookID)
	}

	// Get or create a new user based on facebook ID and other details
	user, err = s.GetAccountsService().GetOrCreateFacebookUser(
		authenticatedAccount,
		facebookID,
		&accounts.UserRequest{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check that the same account is being used
	if authenticatedAccount.ID != user.Account.ID {
		response.UnauthorizedError(w, ErrAccountMismatch.Error())
		return
	}

	// Log in the user
	accessToken, refreshToken, err := s.GetAccountsService().GetOauthService().Login(
		user.Account.OauthClient,
		user.OauthUser,
		scope,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON access token to the response
	accessTokenRespone := &oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    s.cnf.Oauth.AccessTokenLifetime,
		TokenType:    oauth.TokenType,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.Token,
	}
	response.WriteJSON(w, accessTokenRespone, 200)
}
