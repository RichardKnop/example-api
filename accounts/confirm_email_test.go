package accounts_test

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/response"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmEmailFailsWithoutAccountAuthentication() {
	// Prepare a request
	bogusUUID := uuid.New()
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/confirmations/%s", bogusUUID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestConfirmEmailReferenceNotFound() {
	bogusUUID := uuid.New()
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/confirmations/%s", bogusUUID),
		nil,
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

	expectedCode := http.StatusNotFound
	expectedMessage := accounts.ErrConfirmationNotFound.Error()
	response.TestResponseForError(suite.T(), w, expectedMessage, expectedCode)
}

func (suite *AccountsTestSuite) TestConfirmEmail() {
	var (
		testOauthUser    *oauth.User
		testUser         *accounts.User
		testConfirmation *accounts.Confirmation
		err              error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		"harold@finch",
		"test_password",
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser = accounts.NewUser(
		suite.accounts[0],
		testOauthUser,
		suite.userRole,
		"", // facebook ID
		"Harold",
		"Finch",
		"",    // picture
		false, // confirmed
	)
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	testUser.Role = suite.userRole

	// Insert a test confirmation
	testConfirmation = accounts.NewConfirmation(testUser)
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Failed to insert a test confirmation")
	testConfirmation.User = testUser

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/accounts/confirmations/%s", testConfirmation.Reference),
		nil,
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

	// Check empty response
	response.TestEmptyResponse(suite.T(), w)

	// Fetch the updated user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the confirmation
	confirmation := new(accounts.Confirmation)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(confirmation).RecordNotFound())

	// And correct data was saved
	assert.True(suite.T(), user.Confirmed)
}
