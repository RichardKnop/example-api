package accounts_test

import (
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

func (suite *AccountsTestSuite) TestGetUserRequiresUserAuthentication() {
	// Prepare a request
	r, err := http.NewRequest("GET", "http://1.2.3.4/v1/users/12345", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AccountsTestSuite) TestGetUserFailsWithoutPermission() {
	testutil.TestGetErrorExpectedResponse(
		suite.T(),
		suite.router,
		fmt.Sprintf("http://1.2.3.4/v1/users/%d", suite.users[2].ID),
		"get_user",
		"test_user_token",
		accounts.ErrGetUserPermission.Error(),
		http.StatusForbidden,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestGetUser() {
	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/users/%d", suite.users[1].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_user", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Fetch the user
	user := new(accounts.User)
	notFound := accounts.UserPreload(suite.db).First(user, suite.users[1].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

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
		Email:     user.OauthUser.Username,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Role:      roles.User,
		Confirmed: user.Confirmed,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
