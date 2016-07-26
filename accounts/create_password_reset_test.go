package accounts_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/response"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestCreatePasswordResetRequiresAccountAuthentication() {
	// Prepare a request
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/accounts/password-reset",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestCreatePasswordReset() {
	// Prepare a request
	payload, err := json.Marshal(&accounts.PasswordResetRequest{
		Email: "test@user",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/accounts/password-reset",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			b64.StdEncoding.EncodeToString([]byte("test_client_1:test_secret")),
		),
	)

	// Mock password reset email
	suite.mockPasswordResetEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_password_reset", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check empty response
	response.TestEmptyResponse(suite.T(), w)

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Wait for the email goroutine to finish
	<-time.After(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AccountsTestSuite) TestCreatePasswordResetSecondTime() {
	// Insert a test password reset
	testPasswordReset := accounts.NewPasswordReset(suite.users[1])
	err := suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Failed to insert a test password reset")
	testPasswordReset.User = suite.users[1]

	// Prepare a request
	payload, err := json.Marshal(&accounts.PasswordResetRequest{
		Email: suite.users[1].OauthUser.Username},
	)
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/accounts/password-reset",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			b64.StdEncoding.EncodeToString([]byte("test_client_1:test_secret")),
		),
	)

	// Mock password reset email
	suite.mockPasswordResetEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_password_reset", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check empty response
	response.TestEmptyResponse(suite.T(), w)

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the last password reset
	passwordReset := new(accounts.PasswordReset)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(passwordReset).RecordNotFound())

	// And correct data was saved
	assert.NotEqual(suite.T(), testPasswordReset.ID, passwordReset.ID)
	assert.Equal(suite.T(), testPasswordReset.User.ID, passwordReset.User.ID)

	// Wait for the email goroutine to finish
	<-time.After(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
