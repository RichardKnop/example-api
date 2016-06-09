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

	"github.com/RichardKnop/recall/oauth"
	pass "github.com/RichardKnop/recall/password"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmPasswordResetFailsWithoutAccountAuthentication() {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()

	suite.service.confirmPasswordResetHandler(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestConfirmPasswordResetReferenceNotFound() {
	bogusUUID := uuid.New()
	payload, err := json.Marshal(&ConfirmPasswordResetRequest{
		PasswordRequest: PasswordRequest{Password: "test_password"},
	})
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/password-resets/%s", bogusUUID),
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

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 404, w.Code) {
		log.Print(w.Body.String())
	}

	// Check the response body
	expectedJSON, err := json.Marshal(
		map[string]string{"error": ErrPasswordResetNotFound.Error()})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(
			suite.T(),
			string(expectedJSON),
			strings.TrimRight(w.Body.String(), "\n"),
			"Body should contain JSON detailing the error",
		)
	}
}

func (suite *AccountsTestSuite) TestConfirmPasswordReset() {
	var (
		testOauthUser     *oauth.User
		testUser          *User
		testPasswordReset *PasswordReset
		err               error
	)

	// Insert a test user
	testOauthUser, err = suite.service.oauthService.CreateUser(
		"harold@finch",
		"", // blank password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser = NewUser(
		suite.accounts[0],
		testOauthUser,
		suite.userRole,
		"", // facebook ID
		"Harold",
		"Finch",
		false, // confirmed
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Insert a test password reset
	testPasswordReset = NewPasswordReset(testUser)
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Failed to insert a test password reset")
	testPasswordReset.User = testUser

	// Prepare a request
	payload, err := json.Marshal(&ConfirmPasswordResetRequest{
		PasswordRequest: PasswordRequest{Password: "test_password"},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/password-resets/%s", testPasswordReset.Reference),
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

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 204, w.Code) {
		log.Print(w.Body.String())
	}

	// Fetch the updated user
	user := new(User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Password reset should have been soft deleted
	assert.True(suite.T(), suite.db.First(new(PasswordReset), testPasswordReset.ID).RecordNotFound())

	// And correct data was saved
	assert.Nil(suite.T(), pass.VerifyPassword(user.OauthUser.Password.String, "test_password"))

	// Check the response body
	assert.Equal(
		suite.T(),
		"", // empty string
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
