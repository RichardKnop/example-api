package accounts_test

import (
	"net/http"

	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/test-util"
)

func (suite *AccountsTestSuite) TestListUsersRequiresUserAuthentication() {
	testutil.TestGetErrorExpectedResponse(
		suite.T(),
		suite.router,
		"http://1.2.3.4/v1/users",
		"list_users",
		"", // no access token
		accounts.ErrUserAuthenticationRequired.Error(),
		http.StatusUnauthorized,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestListUsersWithoutPermission() {
	testutil.TestGetErrorExpectedResponse(
		suite.T(),
		suite.router,
		"http://1.2.3.4/v1/users",
		"list_users",
		"test_user_token",
		accounts.ErrListUsersPermission.Error(),
		http.StatusForbidden,
		suite.assertMockExpectations,
	)
}

func (suite *AccountsTestSuite) TestListUsers() {
	var users []*accounts.User
	err := accounts.UserPreload(suite.db).Order("id").Find(&users).Error
	assert.NoError(suite.T(), err, "Fetching test users should not fail")

	userResponses := make([]interface{}, len(users))

	for i, user := range users {
		userResponse, err := accounts.NewUserResponse(user)
		assert.NoError(suite.T(), err, "Creating user response should not fail")
		userResponses[i] = userResponse
	}

	testutil.TestListValidResponse(
		suite.T(),
		suite.router,
		"users",                // path
		"users",                // entity
		"test_superuser_token", // from fixtures
		userResponses,
		suite.assertMockExpectations,
	)
}
