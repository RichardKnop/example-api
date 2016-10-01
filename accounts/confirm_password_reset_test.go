package accounts_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/roles"
	pass "github.com/RichardKnop/example-api/password"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmPasswordResetRequiresAccountAuthentication() {
	bogusUUID := uuid.New()
	testutil.TestPostErrorExpectedResponse(
		suite.T(),
		suite.router,
		fmt.Sprintf("http://1.2.3.4/v1/password-resets/%s", bogusUUID),
		"confirm_password_reset",
		nil,
		"", // no access token
		accounts.ErrAccountAuthenticationRequired.Error(),
		http.StatusUnauthorized,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestConfirmPasswordResetNotFound() {
	bogusUUID := uuid.New()
	testutil.TestPostErrorExpectedResponse(
		suite.T(),
		suite.router,
		fmt.Sprintf("http://1.2.3.4/v1/password-resets/%s", bogusUUID),
		"confirm_password_reset",
		strings.NewReader("{}"), //data
		"test_client_token",
		accounts.ErrPasswordResetNotFound.Error(),
		http.StatusNotFound,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestConfirmPasswordReset() {
	var (
		testOauthUser     *oauth.User
		testUser          *accounts.User
		testPasswordReset *accounts.PasswordReset
		err               error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"", // blank password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = accounts.NewUser(
		suite.accounts[0],
		testOauthUser,
		"",    //facebook ID
		false, // confirmed
		&accounts.UserRequest{
			FirstName: "Harold",
			LastName:  "Finch",
		},
	)
	assert.NoError(suite.T(), err, "Failed to create a new user object")
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser

	// Insert a test password reset
	testPasswordReset, err = accounts.NewPasswordReset(
		testUser,
		suite.cnf.AppSpecific.PasswordResetLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new password reset object")
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Failed to insert a test password reset")
	testPasswordReset.User = testUser

	// Prepare a request
	payload, err := json.Marshal(&accounts.ConfirmPasswordResetRequest{
		PasswordRequest: accounts.PasswordRequest{Password: "test_password"},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/password-resets/%s", testPasswordReset.Reference),
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

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "confirm_password_reset", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.PasswordReset)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)

	// Fetch the updated user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Password reset should have been soft deleted
	assert.True(suite.T(), suite.db.First(new(accounts.PasswordReset), testPasswordReset.ID).RecordNotFound())

	// And correct data was saved
	assert.Nil(suite.T(), pass.VerifyPassword(user.OauthUser.Password.String, "test_password"))

	// Check the response
	expected, err := accounts.NewPasswordResetResponse(testPasswordReset)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
