package accounts_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestCreatePasswordResetRequiresClientAuthentication() {
	testutil.TestPostErrorExpectedResponse(
		suite.T(),
		suite.router,
		"http://1.2.3.4/v1/password-resets",
		"create_password_reset",
		strings.NewReader("{}"), // data
		"", // no access token
		accounts.ErrClientAuthenticationRequired.Error(),
		http.StatusUnauthorized,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestCreatePasswordReset() {
	// Prepare a request
	payload, err := json.Marshal(&accounts.PasswordResetRequest{
		Email: "test@user",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/password-resets",
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

	suite.service.WaitForNotifications(1)
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
	suite.db.Model(new(models.PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var countAfter int
	suite.db.Model(new(models.PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created password reset
	passwordReset := new(models.PasswordReset)
	notFound := models.PasswordResetPreload(suite.db).Last(passwordReset).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check the response
	expected, err := accounts.NewPasswordResetResponse(passwordReset)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 201)

	// block until email goroutine has finished
	assert.True(suite.T(), <-suite.service.GetNotifications(), "The email goroutine should have run")

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AccountsTestSuite) TestCreatePasswordResetSecondTime() {
	// Insert a test password reset
	testPasswordReset, err := models.NewPasswordReset(
		suite.users[1],
		suite.cnf.AppSpecific.PasswordResetLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new password reset object")
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Failed to insert a test password reset")
	testPasswordReset.User = suite.users[1]

	// Prepare a request
	payload, err := json.Marshal(&accounts.PasswordResetRequest{
		Email: suite.users[1].OauthUser.Username},
	)
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/password-resets",
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

	suite.service.WaitForNotifications(1)
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
	suite.db.Model(new(models.PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var countAfter int
	suite.db.Model(new(models.PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the created password reset
	passwordReset := new(models.PasswordReset)
	notFound := models.PasswordResetPreload(suite.db).Last(passwordReset).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.NotEqual(suite.T(), testPasswordReset.ID, passwordReset.ID)
	assert.Equal(suite.T(), testPasswordReset.User.ID, passwordReset.User.ID)

	// Check the response
	expected, err := accounts.NewPasswordResetResponse(passwordReset)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 201)

	// block until email goroutine has finished
	assert.True(suite.T(), <-suite.service.GetNotifications(), "The email goroutine should have run")

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
