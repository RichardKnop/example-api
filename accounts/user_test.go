package accounts

import (
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindUserByOauthUserID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByOauthUserID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByOauthUserID(suite.users[1].OauthUser.ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	}
}

func (suite *AccountsTestSuite) TestFindUserByID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	}
}

func (suite *AccountsTestSuite) TestCreateSuperadmin() {
	var (
		user *User
		err  error
	)

	// We try to insert a user with a non unique oauth user
	user, err = suite.service.CreateSuperadmin(
		suite.accounts[0], // account
		"test@superadmin", // email
		"test_password",   // password
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_users.username", err.Error())
	}

	// We try to insert a unique superuser
	user, err = suite.service.CreateSuperadmin(
		suite.accounts[0],  // account
		"test@superadmin2", // email
		"test_password",    // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@superadmin2", user.OauthUser.Username)
	}
}
