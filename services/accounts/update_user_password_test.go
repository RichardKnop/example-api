package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/oauth"
	"github.com/RichardKnop/example-api/services/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/example-api/util/password"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestUpdateUserPasswordFailsWithBadCurrentPassword() {
	var (
		testOauthUser   *models.OauthUser
		testUser        *models.User
		testAccessToken *models.OauthAccessToken
		err             error
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

	// Login the test user
	testAccessToken, _, err = suite.service.GetOauthService().Login(
		suite.clients[0],
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&accounts.UserRequest{
		Password:    "bogus_password",
		NewPassword: "some_new_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the response
	testutil.TestResponseForError(suite.T(), w, oauth.ErrInvalidUserPassword.Error(), 400)
}

func (suite *AccountsTestSuite) TestUpdateUserPasswordFailsWithPaswordlessUser() {
	var (
		testOauthUser   *models.OauthUser
		testUser        *models.User
		testAccessToken *models.OauthAccessToken
		err             error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"", // empty password
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

	// Login the test user
	testAccessToken, _, err = suite.service.GetOauthService().Login(
		suite.clients[0],
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&accounts.UserRequest{
		Password:    "",
		NewPassword: "some_new_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the response
	testutil.TestResponseForError(suite.T(), w, oauth.ErrUserPasswordNotSet.Error(), 400)
}

func (suite *AccountsTestSuite) TestUpdateUserPassword() {
	var (
		testOauthUser   *models.OauthUser
		testUser        *models.User
		testAccessToken *models.OauthAccessToken
		err             error
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

	// Login the test user
	testAccessToken, _, err = suite.service.GetOauthService().Login(
		suite.clients[0],
		testUser.OauthUser,
		"read_write", // scope
	)
	assert.NoError(suite.T(), err, "Failed to login the test user")

	payload, err := json.Marshal(&accounts.UserRequest{
		Password:    "test_password",
		NewPassword: "some_new_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://1.2.3.4/v1/users/%d", testUser.ID),
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testAccessToken.Token))

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "update_user", match.Route.GetName())
	}

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Count after
	var countAfter int
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the password has changed
	assert.Error(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"test_password",
	))
	assert.NoError(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"some_new_password",
	))

	// And the user meta data is unchanged
	assert.Equal(suite.T(), "harold@finch", user.OauthUser.Username)
	assert.Equal(suite.T(), "Harold", user.FirstName.String)
	assert.Equal(suite.T(), "Finch", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	assert.False(suite.T(), user.Confirmed)

	// Check the response
	expected := &accounts.UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/users/%d", user.ID),
				},
			},
		},
		ID:        user.ID,
		Email:     "harold@finch",
		FirstName: "Harold",
		LastName:  "Finch",
		Role:      roles.User,
		Confirmed: false,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
