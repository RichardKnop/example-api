package accounts_test

import (
	"testing"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/example-api/util/pagination"
	"github.com/stretchr/testify/assert"
)

func TestUserGetName(t *testing.T) {
	user := new(models.User)
	assert.Equal(t, "", user.GetName())

	user.FirstName = util.StringOrNull("John")
	user.LastName = util.StringOrNull("Reese")
	assert.Equal(t, "John Reese", user.GetName())
}

func (suite *AccountsTestSuite) TestFindUserByOauthUserID() {
	var (
		user *models.User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByOauthUserID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByOauthUserID(suite.users[1].OauthUser.ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	}
}

func (suite *AccountsTestSuite) TestFindUserByEmail() {
	var (
		user *models.User
		err  error
	)

	// Let's try to find a user by a bogus email
	user, err = suite.service.FindUserByEmail("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Now let's pass a valid email
	user, err = suite.service.FindUserByEmail("test@user")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	}

	// Test case insensitiviness
	user, err = suite.service.FindUserByEmail("TeSt@UsEr")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	}
}

func (suite *AccountsTestSuite) TestFindUserByID() {
	var (
		user *models.User
		err  error
	)

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByID(12345)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByID(suite.users[1].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	}
}

func (suite *AccountsTestSuite) TestFindUserByFacebookID() {
	var (
		user *models.User
		err  error
	)

	// Let's try to find a user by an empty string Facebook ID
	user, err = suite.service.FindUserByFacebookID("")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Let's try to find a user by a bogus ID
	user, err = suite.service.FindUserByFacebookID("bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Now let's pass a valid ID
	user, err = suite.service.FindUserByFacebookID(suite.users[1].FacebookID.String)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user should be returned with preloaded data
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_client_1", user.Account.OauthClient.Key)
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	}
}

func (suite *AccountsTestSuite) TestGetOrCreateFacebookUser() {
	var (
		countBefore, countAfter int
		user                    *models.User
		err                     error
	)

	// Count before
	suite.db.Model(new(models.User)).Count(&countBefore)

	// Let's try passing an existing facebook ID
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0], // account
		"facebook_id_2",   // facebook ID
		&accounts.UserRequest{
			Email:     "test@user",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// Count after
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), "facebook_id_2", user.FacebookID.String)
		assert.Equal(suite.T(), "John", user.FirstName.String)
		assert.Equal(suite.T(), "Reese", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}

	// Count before
	suite.db.Model(new(models.User)).Count(&countBefore)

	// Let's try passing an existing email
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0],   // account
		"finch_facebook_id", // facebook ID
		&accounts.UserRequest{
			Email:     "test@user",
			FirstName: "Harold",
			LastName:  "Finch",
		},
	)

	// Count after
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
		assert.Equal(suite.T(), "finch_facebook_id", user.FacebookID.String)
		assert.Equal(suite.T(), "Harold", user.FirstName.String)
		assert.Equal(suite.T(), "Finch", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}

	// Count before
	suite.db.Model(new(models.User)).Count(&countBefore)

	// We pass new facebook ID and new email
	user, err = suite.service.GetOrCreateFacebookUser(
		suite.accounts[0],   // account
		"reese_facebook_id", // facebook ID
		&accounts.UserRequest{
			Email:     "john@reese",
			FirstName: "John",
			LastName:  "Reese",
		},
	)

	// Count after
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
		assert.Equal(suite.T(), "john@reese", user.OauthUser.Username)
		assert.Equal(suite.T(), "reese_facebook_id", user.FacebookID.String)
		assert.Equal(suite.T(), "John", user.FirstName.String)
		assert.Equal(suite.T(), "Reese", user.LastName.String)
		assert.True(suite.T(), user.Confirmed)
	}
}

func (suite *AccountsTestSuite) TestCreateSuperuser() {
	var (
		user *models.User
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
		assert.Equal(suite.T(), oauth.ErrUsernameTaken, err)
	}

	// We try to insert a unique superuser
	user, err = suite.service.CreateSuperuser(
		suite.accounts[0],   // account
		"test@newsuperuser", // email
		"test_password",     // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
		assert.Equal(suite.T(), "test@newsuperuser", user.OauthUser.Username)
		assert.False(suite.T(), user.FirstName.Valid)
		assert.False(suite.T(), user.LastName.Valid)
		assert.True(suite.T(), user.Confirmed)
	}
}

func (suite *AccountsTestSuite) TestPaginatedUsersCount() {
	var (
		count int
		err   error
	)

	count, err = suite.service.PaginatedUsersCount()
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, count)
	}
}

func (suite *AccountsTestSuite) TestFindPaginatedUsers() {
	var (
		users []*models.User
		err   error
	)

	// This should return all users
	users, err = suite.service.FindPaginatedUsers(
		0,                   // offset
		25,                  // limit
		map[string]string{}, // sorts
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, len(users))
		assert.Equal(suite.T(), suite.users[0].ID, users[0].ID)
		assert.Equal(suite.T(), suite.users[1].ID, users[1].ID)
		assert.Equal(suite.T(), suite.users[2].ID, users[2].ID)
	}

	// This should return all users ordered by ID desc
	users, err = suite.service.FindPaginatedUsers(
		0,  // offset
		25, // limit
		map[string]string{"id": pagination.Descending}, // sorts
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, len(users))
		assert.Equal(suite.T(), suite.users[2].ID, users[0].ID)
		assert.Equal(suite.T(), suite.users[1].ID, users[1].ID)
		assert.Equal(suite.T(), suite.users[0].ID, users[2].ID)
	}

	// Test offset
	users, err = suite.service.FindPaginatedUsers(
		2,                   // offset
		25,                  // limit
		map[string]string{}, // sorts
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(users))
		assert.Equal(suite.T(), suite.users[2].ID, users[0].ID)
	}

	// Test limit
	users, err = suite.service.FindPaginatedUsers(
		1,                   // offset
		1,                   // limit
		map[string]string{}, // sorts
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(users))
		assert.Equal(suite.T(), suite.users[1].ID, users[0].ID)
	}
}
