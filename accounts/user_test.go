package accounts

import (
	"testing"

	"github.com/RichardKnop/recall/util"
	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/stretchr/testify/assert"
)

func TestUserGetName(t *testing.T) {
	user := &User{}
	assert.Equal(t, "", user.GetName())

	user.FirstName = util.StringOrNull("John")
	user.LastName = util.StringOrNull("Reese")
	assert.Equal(t, "John Reese", user.GetName())
}

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
		assert.Equal(suite.T(), ErrUserNotFound, err)
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

func (suite *AccountsTestSuite) TestFindUserByEmail() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by a bogus email
	user, err = suite.service.FindUserByEmail("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid email
	user, err = suite.service.FindUserByEmail("test@user")

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
		assert.Equal(suite.T(), ErrUserNotFound, err)
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

func (suite *AccountsTestSuite) TestFindUserByFacebookID() {
	var (
		user *User
		err  error
	)

	// Let's try to find a user by an empty string Facebook ID
	user, err = suite.service.FindUserByFacebookID("")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByFacebookID("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByFacebookID(suite.users[1].FacebookID.String)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.Role.Name)
	}
}

func (suite *AccountsTestSuite) TestCreateFacebookUser() {
	var (
		user *User
		err  error
	)

	// We try to create a facebook user with a non unique email
	user, err = suite.service.CreateFacebookUser(
		suite.accounts[0], // account
		"facebook_id_3",   // email
		&UserRequest{
			Email:     "test@user",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_users.username", err.Error())
	}

	// We try to create a facebook user with a non unique facebook ID
	user, err = suite.service.CreateFacebookUser(
		suite.accounts[0], // account
		"facebook_id_2",   // email
		&UserRequest{
			Email:     "test@user2",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: account_users.facebook_id", err.Error())
	}

	// We try to get or create a new facebook user
	user, err = suite.service.CreateFacebookUser(
		suite.accounts[0], // account
		"facebook_id_3",   // email
		&UserRequest{
			Email:     "test@user2",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user2", user.OauthUser.Username)
		assert.Equal(suite.T(), "facebook_id_3", user.FacebookID.String)
		assert.Equal(suite.T(), "John", user.FirstName.String)
		assert.Equal(suite.T(), "Reese", user.LastName.String)
	}
}

func (suite *AccountsTestSuite) TestCreateSuperuser() {
	var (
		user *User
		err  error
	)

	// We try to insert a user with a non unique oauth user
	user, err = suite.service.CreateSuperuser(
		suite.accounts[0], // account
		"test@superuser",  // email
		"test_password",   // password
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_users.username", err.Error())
	}

	// We try to insert a unique superuser
	user, err = suite.service.CreateSuperuser(
		suite.accounts[0], // account
		"test@superuser2", // email
		"test_password",   // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@superuser2", user.OauthUser.Username)
		assert.False(suite.T(), user.FacebookID.Valid)
		assert.False(suite.T(), user.FirstName.Valid)
		assert.False(suite.T(), user.LastName.Valid)
	}
}
