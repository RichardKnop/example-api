package accounts

import (
	"net/http"
	"testing"

	"github.com/RichardKnop/recall/util"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthenticatedAccount(t *testing.T) {
	var (
		account *Account
		err     error
	)

	// A test request
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	account, err = GetAuthenticatedAccount(r)

	// Account object should be nil
	assert.Nil(t, account)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, errAccountAuthenticationRequired, err)
	}

	// Set a context value of an invalid type
	context.Set(r, authenticatedAccountKey, "bogus")

	account, err = GetAuthenticatedAccount(r)

	// Account object should be nil
	assert.Nil(t, account)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, errAccountAuthenticationRequired, err)
	}

	// Set a valid context value
	context.Set(r, authenticatedAccountKey, &Account{Name: "Test Account"})

	account, err = GetAuthenticatedAccount(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct account object should be returned
	if assert.NotNil(t, account) {
		assert.Equal(t, "Test Account", account.Name)
	}
}

func TestGetAuthenticatedUser(t *testing.T) {
	var (
		user *User
		err  error
	)

	// A test request
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	user, err = GetAuthenticatedUser(r)

	// User object should be nil
	assert.Nil(t, user)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, errUserAuthenticationRequired, err)
	}

	// Set a context value of an invalid type
	context.Set(r, authenticatedUserKey, "bogus")

	user, err = GetAuthenticatedUser(r)

	// User object should be nil
	assert.Nil(t, user)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, errUserAuthenticationRequired, err)
	}

	// Set a valid context value
	context.Set(r, authenticatedUserKey, &User{FirstName: util.StringOrNull("John Reese")})

	user, err = GetAuthenticatedUser(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct user object should be returned
	if assert.NotNil(t, user) {
		assert.Equal(t, "John Reese", user.FirstName.String)
	}
}
