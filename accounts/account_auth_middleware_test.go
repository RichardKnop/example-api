package accounts_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/recall/accounts"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestAccountAuthMiddleware() {
	var (
		r                    *http.Request
		w                    *httptest.ResponseRecorder
		next                 http.HandlerFunc
		authenticatedAccount *accounts.Account
		err                  error
	)

	middleware := accounts.NewAccountAuthMiddleware(suite.service)

	// Send a request without basic auth through the middleware
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
			accounts.ErrAccountAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedAccount, err = accounts.GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), authenticatedAccount)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrAccountAuthenticationRequired, err)
	}

	// Send a request with incorrect basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("bogus", "bogus")
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
			accounts.ErrAccountAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedAccount, err = accounts.GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), authenticatedAccount)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrAccountAuthenticationRequired, err)
	}

	// Send a request with correct basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variable has been set
	authenticatedAccount, err = accounts.GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), err)
	if assert.NotNil(suite.T(), authenticatedAccount) {
		assert.Equal(suite.T(), "Test Account 1", authenticatedAccount.Name)
		assert.Equal(suite.T(), "test_client_1", authenticatedAccount.OauthClient.Key)
	}
}
