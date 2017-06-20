package accounts_test

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/services/oauth"
	"github.com/RichardKnop/example-api/services/oauth/roles"
	"github.com/RichardKnop/example-api/services/oauth/tokentypes"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestConfirmEmailRequiresClientAuthentication() {
	bogusUUID := uuid.New()
	testutil.TestGetErrorExpectedResponse(
		suite.T(),
		suite.router,
		fmt.Sprintf("http://1.2.3.4/v1/confirmations/%s", bogusUUID),
		"confirm_email",
		"", // no access token
		accounts.ErrClientAuthenticationRequired.Error(),
		http.StatusUnauthorized,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestConfirmEmailNotFound() {
	bogusUUID := uuid.New()
	testutil.TestGetErrorExpectedResponse(
		suite.T(),
		suite.router,
		fmt.Sprintf("http://1.2.3.4/v1/confirmations/%s", bogusUUID),
		"confirm_email",
		"test_client_token",
		accounts.ErrConfirmationNotFound.Error(),
		http.StatusNotFound,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestConfirmEmail() {
	var (
		testOauthUser    *models.OauthUser
		testUser         *models.User
		testConfirmation *models.Confirmation
		err              error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"test_password",
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = models.NewUser(
		suite.clients[0],
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
	testUser.OauthClient = suite.clients[0]
	testUser.OauthUser = testOauthUser

	// Insert a test confirmation
	testConfirmation, err = models.NewConfirmation(
		testUser,
		suite.cnf.AppSpecific.ConfirmationLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new confirmation object")
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Failed to insert a test confirmation")
	testConfirmation.User = testUser

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/confirmations/%s",
			testConfirmation.Reference,
		),
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

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "confirm_email", match.Route.GetName())
	}

	// Count before
	var (
		countBefore              int
		accessTokensCountBefore  int
		refreshTokensCountBefore int
	)
	suite.db.Model(new(models.Confirmation)).Count(&countBefore)
	suite.db.Model(new(models.OauthAccessToken)).Count(&accessTokensCountBefore)
	suite.db.Model(new(models.OauthRefreshToken)).Count(&refreshTokensCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var (
		countAfter              int
		accessTokensCountAfter  int
		refreshTokensCountAfter int
	)
	suite.db.Model(new(models.Confirmation)).Count(&countAfter)
	suite.db.Model(new(models.OauthAccessToken)).Count(&accessTokensCountAfter)
	suite.db.Model(new(models.OauthRefreshToken)).Count(&refreshTokensCountAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)
	assert.Equal(suite.T(), accessTokensCountBefore, accessTokensCountAfter)
	assert.Equal(suite.T(), refreshTokensCountBefore, refreshTokensCountAfter)

	// Fetch the updated user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Confirmation should have been soft deleteted
	assert.True(suite.T(), suite.db.Last(new(models.Confirmation)).RecordNotFound())

	// And correct data was saved
	assert.True(suite.T(), user.Confirmed)

	// Check the response
	expected, err := accounts.NewConfirmationResponse(testConfirmation)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}

func (suite *AccountsTestSuite) TestConfirmEmailWithAutologinFlag() {
	var (
		testOauthUser    *models.OauthUser
		testUser         *models.User
		testConfirmation *models.Confirmation
		err              error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"test_password",
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = models.NewUser(
		suite.clients[0],
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
	testUser.OauthClient = suite.clients[0]

	// Insert a test confirmation
	testConfirmation, err = models.NewConfirmation(
		testUser,
		suite.cnf.AppSpecific.ConfirmationLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new confirmation object")
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Failed to insert a test confirmation")
	testConfirmation.User = testUser

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://1.2.3.4/v1/confirmations/%s?autologin=true",
			testConfirmation.Reference,
		),
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

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "confirm_email", match.Route.GetName())
	}

	// Count before
	var (
		countBefore              int
		accessTokensCountBefore  int
		refreshTokensCountBefore int
	)
	suite.db.Model(new(models.Confirmation)).Count(&countBefore)
	suite.db.Model(new(models.OauthAccessToken)).Count(&accessTokensCountBefore)
	suite.db.Model(new(models.OauthRefreshToken)).Count(&refreshTokensCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Count after
	var (
		countAfter              int
		accessTokensCountAfter  int
		refreshTokensCountAfter int
	)
	suite.db.Model(new(models.Confirmation)).Count(&countAfter)
	suite.db.Model(new(models.OauthAccessToken)).Count(&accessTokensCountAfter)
	suite.db.Model(new(models.OauthRefreshToken)).Count(&refreshTokensCountAfter)
	assert.Equal(suite.T(), countBefore-1, countAfter)
	assert.Equal(suite.T(), accessTokensCountBefore+1, accessTokensCountAfter)
	assert.Equal(suite.T(), refreshTokensCountBefore+1, refreshTokensCountAfter)

	// Fetch the updated user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Confirmation should have been soft deleteted
	assert.True(suite.T(), suite.db.Last(new(models.Confirmation)).RecordNotFound())

	// And correct data was saved
	assert.True(suite.T(), user.Confirmed)

	// Fetch login data
	accessToken, refreshToken := new(models.OauthAccessToken), new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		Last(accessToken).RecordNotFound())
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		Last(refreshToken).RecordNotFound())

	// Check the response
	expectedEmbeddedTokenResponse, err := oauth.NewAccessTokenResponse(
		accessToken,
		refreshToken,
		suite.cnf.Oauth.AccessTokenLifetime,
		tokentypes.Bearer,
	)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	expected, err := accounts.NewConfirmationResponse(testConfirmation)
	expected.SetEmbedded(
		"access-token",
		jsonhal.Embedded(expectedEmbeddedTokenResponse),
	)
	assert.NoError(suite.T(), err, "Failed to create expected response object")
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
