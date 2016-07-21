package facebook

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/response"
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

	// Initialise variables
	var (
		profile UserProfile
		user    *accounts.User
	)

	// Decode the response to struct
	if err := resp.Decode(&profile); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if profile.Email == nil {
		// There is an edge case where Facebook does not return a valid email
		// User could have registered with a phone number or have an unconfirmed
		// email address. In such rare case, default to {facebook_id}@facebook.com
		edgeCaseEmail := fmt.Sprintf("%s@facebook.com", profile.ID)
		profile.Email = &edgeCaseEmail
	}

	// Build user request object
	userRequest := &accounts.UserRequest{
		Email:     *profile.Email,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Picture:   profile.GetPictureURL(),
	}

	// Get or create a new user based on facebook ID and other details
	user, err = s.GetAccountsService().GetOrCreateFacebookUser(
		authenticatedAccount,
		profile.ID,
		userRequest,
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
