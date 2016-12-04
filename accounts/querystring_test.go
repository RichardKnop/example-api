package accounts_test

import (
	"fmt"
	"net/http"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestGetUserFromQueryString() {
	var (
		r    *http.Request
		user *models.User
		err  error
	)

	// Let's try with a bogus user ID, should fail
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foobar?user_id=9999", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	user, err = suite.service.GetUserFromQueryString(r)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrUserNotFound, err)
	}

	// Let's try with a valid user ID
	r, err = http.NewRequest(
		"GET",
		fmt.Sprintf("http://1.2.3.4/v1/foobar?user_id=%d", suite.users[1].ID),
		nil,
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	user, err = suite.service.GetUserFromQueryString(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	}
}
