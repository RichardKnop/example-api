package accounts

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestAccountAuthMiddleware() {
	var (
		r                    *http.Request
		w                    *httptest.ResponseRecorder
		next                 http.HandlerFunc
		authenticatedAccount *Account
		err                  error
	)

	middleware := NewAccountAuthMiddleware(suite.service)

	// Send a request without basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
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
			errAccountAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedAccount, err = GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), authenticatedAccount)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAccountAuthenticationRequired, err)
	}

	// Send a request with incorrect basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
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
			errAccountAuthenticationRequired,
		),
		strings.TrimSpace(w.Body.String()),
	)

	// Check the context variable has not been set
	authenticatedAccount, err = GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), authenticatedAccount)
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAccountAuthenticationRequired, err)
	}

	// Send a request with correct basic auth through the middleware
	r, err = http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth("test_client_1", "test_secret")
	w = httptest.NewRecorder()
	next = func(w http.ResponseWriter, r *http.Request) {}
	middleware.ServeHTTP(w, r, next)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the context variable has been set
	authenticatedAccount, err = GetAuthenticatedAccount(r)
	assert.Nil(suite.T(), err)
	if assert.NotNil(suite.T(), authenticatedAccount) {
		assert.Equal(suite.T(), "Test Account 1", authenticatedAccount.Name)
		assert.Equal(suite.T(), "test_client_1", authenticatedAccount.OauthClient.Key)
	}
}
