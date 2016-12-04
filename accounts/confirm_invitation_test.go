package accounts_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	pass "github.com/RichardKnop/example-api/util/password"
	"github.com/RichardKnop/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmInvitationFailsWithoutAccountAuthentication() {
	// Prepare a request
	bogusUUID := uuid.New()
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/invitations/%s", bogusUUID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestConfirmInvitationReferenceNotFound() {
	// Prepare a request
	bogusUUID := uuid.New()
	payload, err := json.Marshal(&accounts.ConfirmInvitationRequest{
		PasswordRequest: accounts.PasswordRequest{Password: "test_password"},
	})
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/invitations/%s", bogusUUID),
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

	// Check the response
	testutil.TestResponseForError(suite.T(), w, accounts.ErrInvitationNotFound.Error(), 404)
}

func (suite *AccountsTestSuite) TestConfirmInvitation() {
	var (
		testOauthUser  *models.OauthUser
		testUser       *models.User
		testInvitation *models.Invitation
		err            error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"", // blank password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = models.NewUser(
		suite.accounts[0],
		testOauthUser,
		"", //facebook ID
		"Harold",
		"Finch",
		"",    // picture
		false, // confirmed
	)
	assert.NoError(suite.T(), err, "Failed to create a new user object")
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser

	// Insert a test invitation
	testInvitation, err = models.NewInvitation(
		testUser,
		suite.users[0],
		suite.cnf.AppSpecific.InvitationLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new invitation object")
	err = suite.db.Create(testInvitation).Error
	assert.NoError(suite.T(), err, "Failed to insert a test invitation")
	testInvitation.InvitedUser = testUser
	testInvitation.InvitedByUser = suite.users[0]

	// Prepare a request
	payload, err := json.Marshal(&accounts.ConfirmInvitationRequest{
		PasswordRequest: accounts.PasswordRequest{Password: "test_password"},
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://1.2.3.4/v1/invitations/%s", testInvitation.Reference),
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
		assert.Equal(suite.T(), "confirm_invitation", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(models.Invitation)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var countAfter int
	suite.db.Model(new(models.Invitation)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)

	// Fetch the updated user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Invitation should have been soft deleted
	assert.True(suite.T(), suite.db.First(new(models.Invitation), testInvitation.ID).RecordNotFound())

	// And correct data was saved
	assert.Nil(suite.T(), pass.VerifyPassword(user.OauthUser.Password.String, "test_password"))

	// Check the response
	expected, err := accounts.NewInvitationResponse(testInvitation)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
