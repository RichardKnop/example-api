package accounts_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestUserAuthMiddleware() {
	var (
		r                 *http.Request
		w                 *httptest.ResponseRecorder
		next              http.HandlerFunc
		authenticatedUser *accounts.User
		err               error
	)

	middleware := accounts.NewUserAuthMiddleware(suite.service)

	// Send a request without a bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf(
			"{\"error\":\"%s\"}",
			accounts.ErrUserAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)
	}

	// Send a request with empty bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer ")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf(
			"{\"error\":\"%s\"}",
			accounts.ErrUserAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)
	}

	// Send a request with incorrect bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer bogus")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf(
			"{\"error\":\"%s\"}",
			accounts.ErrUserAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)
	}

	// Send a request with client bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_client_token")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf(
			"{\"error\":\"%s\"}",
			accounts.ErrUserAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), authenticatedUser)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserAuthenticationRequired, err)
	}

	// Send a request with correct bearer token through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variable has been set
	authenticatedUser, err = accounts.GetAuthenticatedUser(r)
	assert.Nil(suite.T(), err)
	if assert.NotNil(suite.T(), authenticatedUser) {
		assert.Equal(suite.T(), "test@user", authenticatedUser.OauthUser.Username)
	}
}
