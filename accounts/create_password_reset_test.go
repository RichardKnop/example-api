package accounts

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestCreatePasswordReset() {
	// Prepare a request
	payload, err := json.Marshal(&PasswordResetRequest{"test@user"})
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

	// Mock confirmation email
	suite.mockPasswordResetEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_password_reset", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Sleep for the email goroutine to finish
	time.Sleep(10 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created password reset
	passwordReset := new(PasswordReset)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(passwordReset).RecordNotFound())

	// And correct data was saved
	assert.Equal(suite.T(), "test@user", passwordReset.User.OauthUser.Username)
	assert.True(suite.T(), passwordReset.EmailSent)
	assert.True(suite.T(), passwordReset.EmailSentAt.Valid)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}

func (suite *AccountsTestSuite) TestCreatePasswordResetSecondTime() {
	// Insert a test password reset
	testPasswordReset := NewPasswordReset(suite.users[1])
	err := suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Failed to insert a test password reset")
	testPasswordReset.User = suite.users[1]

	// Prepare a request
	payload, err := json.Marshal(&PasswordResetRequest{suite.users[1].OauthUser.Username})
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

	// Mock confirmation email
	suite.mockPasswordResetEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_password_reset", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the last password reset
	passwordReset := new(PasswordReset)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(passwordReset).RecordNotFound())

	// And correct data was saved
	assert.NotEqual(suite.T(), testPasswordReset.ID, passwordReset.ID)
	assert.Equal(suite.T(), testPasswordReset.User.ID, passwordReset.User.ID)

	// Email should not have been sent yet
	assert.False(suite.T(), passwordReset.EmailSent)
	assert.False(suite.T(), passwordReset.EmailSentAt.Valid)

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)

	// Sleep for the email goroutine to finish
	time.Sleep(10 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Refresh the password reset
	passwordReset = new(PasswordReset)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(passwordReset).RecordNotFound())

	// Email should have been sent
	assert.True(suite.T(), passwordReset.EmailSent)
	assert.True(suite.T(), passwordReset.EmailSentAt.Valid)
}
