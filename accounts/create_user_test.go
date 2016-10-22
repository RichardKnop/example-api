package accounts_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestCreateUserRequiresAccountAuthentication() {
	// Prepare a request
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/users",
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated account")
}

func (suite *AccountsTestSuite) TestCreateUserOnlyRequiredFields() {
	// Prepare a request
	payload, err := json.Marshal(&accounts.UserRequest{
		Email:    "test@newuser",
		Password: "test_password",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/users",
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
	// Mock confirmation email
	suite.mockConfirmationEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_user", match.Route.GetName())
	}

	// Count before
	var (
		countBefore              int
		confirmationsCountBefore int
	)
	suite.db.Model(new(accounts.User)).Count(&countBefore)
	suite.db.Model(new(accounts.Confirmation)).Count(&confirmationsCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)
	// block until email goroutine has finished
	assert.True(suite.T(), <-suite.service.GetNotifications(), "The email goroutine should have run")

	// Count after
	var (
		countAfter              int
		confirmationsCountAfter int
	)
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	suite.db.Model(new(accounts.Confirmation)).Count(&confirmationsCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), confirmationsCountBefore+1, confirmationsCountAfter)

	// Fetch the created user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created confirmation
	confirmation := new(accounts.Confirmation)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(confirmation).RecordNotFound())

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "test@newuser", user.OauthUser.Username)
	assert.False(suite.T(), user.FirstName.Valid)
	assert.False(suite.T(), user.LastName.Valid)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	assert.False(suite.T(), user.Confirmed)
	assert.Equal(suite.T(), "test@newuser", confirmation.User.OauthUser.Username)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/users/%d", user.ID),
		w.Header().Get("Location"),
	)

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
		Email:     "test@newuser",
		Role:      roles.User,
		Confirmed: false,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 201)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}

func (suite *AccountsTestSuite) TestCreateUserWithOptionalFields() {
	// Prepare a request
	payload, err := json.Marshal(&accounts.UserRequest{
		Email:     "test@newuser",
		Password:  "test_password",
		FirstName: "John",
		LastName:  "Reese",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/users",
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
	// Mock confirmation email
	suite.mockConfirmationEmail()

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "create_user", match.Route.GetName())
	}

	// Count before
	var (
		countBefore              int
		confirmationsCountBefore int
	)
	suite.db.Model(new(accounts.User)).Count(&countBefore)
	suite.db.Model(new(accounts.Confirmation)).Count(&confirmationsCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)
	// block until email goroutine has finished
	assert.True(suite.T(), <-suite.service.GetNotifications(), "The email goroutine should have run")

	// Count after
	var (
		countAfter              int
		confirmationsCountAfter int
	)
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	suite.db.Model(new(accounts.Confirmation)).Count(&confirmationsCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), confirmationsCountBefore+1, confirmationsCountAfter)

	// Fetch the created user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Fetch the created confirmation
	confirmation := new(accounts.Confirmation)
	assert.False(suite.T(), suite.db.Preload("User.OauthUser").
		Last(confirmation).RecordNotFound())

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "test@newuser", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	assert.False(suite.T(), user.Confirmed)
	assert.Equal(suite.T(), "test@newuser", confirmation.User.OauthUser.Username)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/users/%d", user.ID),
		w.Header().Get("Location"),
	)

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
		Email:     "test@newuser",
		Role:      roles.User,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Confirmed: false,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 201)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
