package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/models"
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
func (s *Service) GetUserCredentialsFromToken(token string) (*models.User, error) {
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

// GetClientCredentialsFromBaseAuth does base auth and returns client object
func (s *Service) GetClientCredentialsFromBaseAuth(r *http.Request) (*models.OauthClient, error) {
	// Try to get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		return nil, ErrBaseAuthRequired
	}

	// Authenticate the client
	return s.GetOauthService().AuthClient(clientID, secret)
}

// GetClientCredentialsFromToken ...
func (s *Service) GetClientCredentialsFromToken(token string) (*models.OauthClient, error) {
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

	return oauthAccessToken.Client, nil
}

// GetMixedCredentialsFromToken ...
func (s *Service) GetMixedCredentialsFromToken(token string) (*models.OauthClient, *models.User, error) {
	// Authenticate
	oauthAccessToken, err := s.GetOauthService().Authenticate(token)
	if err != nil {
		return nil, nil, err
	}

	// This should never happen
	if !oauthAccessToken.ClientID.Valid {
		return nil, nil, ErrAccessTokenWithoutClientID
	}

	return s.findClientOrUserFromToken(oauthAccessToken)
}

// getClientCredentialsFromRequest is the common code used to parse client
// credentials in the request object. It will return a client object based on
// either base auth client ID and secret or a client only bearer token
func getClientCredentialsFromRequest(r *http.Request, service ServiceInterface) (*models.OauthClient, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		return service.GetClientCredentialsFromBaseAuth(r)
	}
	return service.GetClientCredentialsFromToken(string(token))
}

// getUserCredentialsFromRequest is the common code used to parse user
// credentials in the request object. It will return a user object based on
// a user specific bearer token
func getUserCredentialsFromRequest(r *http.Request, service ServiceInterface) (*models.User, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		return nil, err
	}
	return service.GetUserCredentialsFromToken(string(token))
}

// getMixedCredentialsFromRequest is the common code used to parse client or
// user credentials in the request object. It will return a client object and/or
// a user based on either base auth or bearer token
func getMixedCredentialsFromRequest(r *http.Request, service ServiceInterface) (*models.OauthClient, *models.User, error) {
	token, err := util.ParseBearerToken(r)
	if err != nil {
		client, err := service.GetClientCredentialsFromBaseAuth(r)
		if err != nil {
			return nil, nil, err
		}
		return client, nil, nil
	}
	return service.GetMixedCredentialsFromToken(string(token))
}

func (s *Service) findClientOrUserFromToken(oauthAccessToken *models.OauthAccessToken) (*models.OauthClient, *models.User, error) {
	// This is a user specific access token
	if oauthAccessToken.UserID.Valid {
		user, err := s.FindUserByOauthUserID(uint(oauthAccessToken.UserID.Int64))
		if err != nil {
			return nil, nil, err
		}
		return nil, user, nil
	}

	return oauthAccessToken.Client, nil, nil
}
