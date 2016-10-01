package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/password"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestUpdateUserRequiresUserAuthentication() {
	testutil.TestPutErrorExpectedResponse(
		suite.T(),
		suite.router,
		"http://1.2.3.4/v1/users/123",
		"update_user",
		nil,
		"", // no access token
		accounts.ErrUserAuthenticationRequired.Error(),
		http.StatusUnauthorized,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestUpdateUser() {
	testUser, testAccessToken, err := suite.insertTestUser(
		"harold@finch", "test_password", "Harold", "Finch")
	assert.NoError(suite.T(), err, "Failed to insert a test user")

	payload, err := json.Marshal(&accounts.UserRequest{
		FirstName: "John",
		LastName:  "Reese",
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
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check that the password has NOT changed
	assert.NoError(suite.T(), password.VerifyPassword(
		user.OauthUser.Password.String,
		"test_password",
	))

	// And correct data was saved
	assert.Equal(suite.T(), "harold@finch", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
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
		FirstName: "John",
		LastName:  "Reese",
		Role:      roles.User,
		Confirmed: false,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
