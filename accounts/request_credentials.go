package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
)

var (
	// ErrBaseAuthRequired ...
	ErrBaseAuthRequired = errors.New("Base auth required")
	// ErrUserAccessTokenRequired ...
	ErrUserAccessTokenRequired = errors.New("Requires user specific access token")
	// ErrClientAccessTokenRequired ...
	ErrClientAccessTokenRequired = errors.New("Requires client only access token")
	// ErrAccessTokenWithoutClientID ...
	ErrAccessTokenWithoutClientID = errors.New("Access token has no client ID")
)

// GetUserCredentialsFromToken ...
func (s *Service) GetUserCredentialsFromToken(token string) (*User, error) {
	// Authenticate
	oauthAccessToken, err := s.GetOauthService().Authenticate(token)
	if err != nil {
		return nil, err
	}

	// This should never happen
	if !oauthAccessToken.UserID.Valid {
		return nil, ErrUserAccessTokenRequired
	}

	return s.FindUserByOauthUserID(uint(oauthAccessToken.UserID.Int64))
}

// GetClientCredentialsFromBaseAuth does base auth and returns account credential
func (s *Service) GetClientCredentialsFromBaseAuth(r *http.Request) (*Account, error) {
	// Try to get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		return nil, ErrBaseAuthRequired
	}

	// Authenticate the client
	oauthClient, err := s.GetOauthService().AuthClient(clientID, secret)
	if err != nil {
		return nil, err
	}

	return s.FindAccountByOauthClientID(oauthClient.ID)
}

// GetClientCredentialsFromToken ...
func (s *Service) GetClientCredentialsFromToken(token string) (*Account, error) {
	// Authenticate
	oauthAccessToken, err := s.GetOauthService().Authenticate(token)
	if err != nil {
		return nil, err
	}

	// This should never happen
	if !oauthAccessToken.ClientID.Valid {
		return nil, ErrAccessTokenWithoutClientID
	}

	// This is a user specific access token
	if oauthAccessToken.UserID.Valid {
		return nil, ErrClientAccessTokenRequired
	}

	return s.FindAccountByOauthClientID(uint(oauthAccessToken.ClientID.Int64))
}

// GetMixedCredentialsFromToken ...
func (s *Service) GetMixedCredentialsFromToken(token string) (*Account, *User, error) {
	// Authenticate
	oauthAccessToken, err := s.GetOauthService().Authenticate(token)
	if err != nil {
		return nil, nil, err
	}

	// This should never happen
	if !oauthAccessToken.ClientID.Valid {
		return nil, nil, ErrAccessTokenWithoutClientID
	}

	return s.findAccountOrUserFromToken(oauthAccessToken)
}

// getClientCredentialsFromRequest is the common code used to parse client
// credentials in the request object. It will return a client object based on
// either base auth client ID and secret or a client only bearer token
func getClientCredentialsFromRequest(r *http.Request, service ServiceInterface) (*Account, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		return service.GetClientCredentialsFromBaseAuth(r)
	}
	account, err := service.GetClientCredentialsFromToken(string(token))
	if err != nil {
		return nil, err
	}
	return account, nil
}

// getUserCredentialsFromRequest is the common code used to parse user
// credentials in the request object. It will return a user object based on
// a user specific bearer token
func getUserCredentialsFromRequest(r *http.Request, service ServiceInterface) (*User, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		return nil, err
	}
	return service.GetUserCredentialsFromToken(string(token))
}

// getMixedCredentialsFromRequest is the common code used to parse client or
// user credentials in the request object. It will return a client object and/or
// a user based on either base auth or bearer token
func getMixedCredentialsFromRequest(r *http.Request, service ServiceInterface) (*Account, *User, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		account, err := service.GetClientCredentialsFromBaseAuth(r)
		if err != nil {
			return nil, nil, err
		}
		return account, nil, nil
	}
	return service.GetMixedCredentialsFromToken(string(token))
}

func (s *Service) findAccountOrUserFromToken(oauthAccessToken *oauth.AccessToken) (*Account, *User, error) {
	// This is a user specific access token
	if oauthAccessToken.UserID.Valid {
		user, err := s.FindUserByOauthUserID(uint(oauthAccessToken.UserID.Int64))
		if err != nil {
			return nil, nil, err
		}
		return nil, user, nil
	}

	// This is a client specific access token
	account, err := s.FindAccountByOauthClientID(uint(oauthAccessToken.ClientID.Int64))
	if err != nil {
		return nil, nil, err
	}

	return account, nil, nil
}
