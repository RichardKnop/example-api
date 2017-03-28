package accounts_test

import (
	"net/http"
	"testing"

	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/util"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthenticatedClient(t *testing.T) {
	var (
		client *models.OauthClient
		err    error
	)

	// A test request
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	client, err = accounts.GetAuthenticatedClient(r)

	// Client object should be nil
	assert.Nil(t, client)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, accounts.ErrClientAuthenticationRequired, err)
	}

	// Set a context value of an invalid type
	context.Set(r, accounts.AuthenticatedClientKey, "bogus")

	client, err = accounts.GetAuthenticatedClient(r)

	// Client object should be nil
	assert.Nil(t, client)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, accounts.ErrClientAuthenticationRequired, err)
	}

	// Set a valid context value
	context.Set(r, accounts.AuthenticatedClientKey, &models.OauthClient{Key: "test_client_1"})

	client, err = accounts.GetAuthenticatedClient(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct client object should be returned
	if assert.NotNil(t, client) {
		assert.Equal(t, "test_client_1", client.Key)
	}
}

func TestGetAuthenticatedUser(t *testing.T) {
	var (
		user *models.User
		err  error
	)

	// A test request
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	user, err = accounts.GetAuthenticatedUser(r)

	// User object should be nil
	assert.Nil(t, user)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, accounts.ErrUserAuthenticationRequired, err)
	}

	// Set a context value of an invalid type
	context.Set(r, accounts.AuthenticatedUserKey, "bogus")

	user, err = accounts.GetAuthenticatedUser(r)

	// User object should be nil
	assert.Nil(t, user)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, accounts.ErrUserAuthenticationRequired, err)
	}

	// Set a valid context value
	context.Set(r, accounts.AuthenticatedUserKey, &models.User{FirstName: util.StringOrNull("John Reese")})

	user, err = accounts.GetAuthenticatedUser(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct user object should be returned
	if assert.NotNil(t, user) {
		assert.Equal(t, "John Reese", user.FirstName.String)
	}
}
