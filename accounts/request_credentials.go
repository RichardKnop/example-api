package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
)

// getCredentialsFromRequest is the common code used to parse credentials in the request object.
// It will return an account object (if available), a client (if available) and any error
func getCredentialsFromRequest(r *http.Request, service ServiceInterface) (*Account, *User, error) {
	token, err := util.ParseBearerToken(r)
	if err == nil {
		return service.getUserCredentials(string(token))
	}
	return service.getClientCredentials(r)
}

func (s *Service) getUserCredentials(token string) (*Account, *User, error) {
	var (
		account          *Account
		user             *User
		oauthAccessToken *oauth.AccessToken
		err              error
	)

	// Authenticate
	oauthAccessToken, err = s.GetOauthService().Authenticate(token)
	if err != nil {
		// For security reasons, return a general error message
		return nil, nil, ErrUserAuthenticationRequired
	}

	if !oauthAccessToken.ClientID.Valid && !oauthAccessToken.UserID.Valid {
		// Needs to have at least one credential in the token
		return nil, nil, ErrUserAuthenticationRequired
	}

	if oauthAccessToken.ClientID.Valid {
		// Fetch the account from the database
		account, err = s.FindAccountByOauthClientID(uint(oauthAccessToken.ClientID.Int64))
		if err != nil {
			// For security reasons, return a general error message
			return nil, nil, ErrUserAuthenticationRequired
		}
	}

	// Access token has no user, this probably means client credentials grant
	if !oauthAccessToken.UserID.Valid {
		return account, nil, nil
	}

	// Fetch the user from the database
	user, err = s.FindUserByOauthUserID(uint(oauthAccessToken.UserID.Int64))
	if err != nil {
		// For security reasons, return a general error message
		return account, nil, ErrUserAuthenticationRequired
	}

	return account, user, nil
}

func (s *Service) getClientCredentials(r *http.Request) (*Account, *User, error) {
	// Try to get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		// For security reasons, return a general error message
		return nil, nil, ErrUserAuthenticationRequired
	}

	// Authenticate the client
	oauthClient, err := s.GetOauthService().AuthClient(clientID, secret)
	if err != nil {
		return nil, nil, ErrUserAuthenticationRequired
	}

	// Fetch the account from the database
	account, err := s.FindAccountByOauthClientID(oauthClient.ID)
	if err != nil {
		return nil, nil, ErrUserAuthenticationRequired
	}

	return account, nil, nil
}
