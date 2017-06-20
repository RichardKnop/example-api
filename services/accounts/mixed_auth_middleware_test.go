package accounts_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestMixedAuthMiddleware() {
	var (
		r                   *http.Request
		w                   *httptest.ResponseRecorder
		next                http.HandlerFunc
		authenticatedClient *models.OauthClient
		authenticatedUser   *models.User
		err                 error
	)

	middleware := accounts.NewMixedAuthMiddleware(suite.service)

	// Send a request without basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		accounts.ErrClientOrUserAuthenticationRequired.Error(),
		401,
	)

	// Check the context variables have not been set
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with incorrect basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("bogus", "bogus")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		accounts.ErrClientOrUserAuthenticationRequired.Error(),
		401,
	)

	// Check the context variables have not been set
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with correct basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variables
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), authenticatedClient)
	assert.Equal(suite.T(), "test_client_1", authenticatedClient.Key)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with correct client access token as bearer
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_client_token")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variables
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), authenticatedClient)
	assert.Equal(suite.T(), "test_client_1", authenticatedClient.Key)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request without a bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		accounts.ErrClientOrUserAuthenticationRequired.Error(),
		401,
	)

	// Check the context variables have not been set
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with empty bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer ")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		accounts.ErrClientOrUserAuthenticationRequired.Error(),
		401,
	)

	// Check the context variables have not been set
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with incorrect bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer bogus")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		accounts.ErrClientOrUserAuthenticationRequired.Error(),
		401,
	)

	// Check the context variables have not been set
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)

	// Send a request with correct bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variables
	authenticatedClient, err = accounts.GetAuthenticatedClient(r)
	assert.Nil(suite.T(), authenticatedClient)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), accounts.ErrClientAuthenticationRequired, err)
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), authenticatedUser)
	assert.Equal(suite.T(), "test@user", authenticatedUser.OauthUser.Username)
}
