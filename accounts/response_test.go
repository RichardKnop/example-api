package accounts_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/oauth"
)

func TestNewUserResponse(t *testing.T) {
	testUser := &accounts.User{
		OauthUser: new(oauth.User),
	}

	// Create user response
	response, err := accounts.NewUserResponse(testUser)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, fmt.Sprintf("/v1/users/%d", testUser.ID), selfLink.Href)
	}
}

func TestNewConfirmationResponse(t *testing.T) {
	testUser := &accounts.User{
		OauthUser: new(oauth.User),
	}
	testConfirmation := &accounts.Confirmation{
		User: testUser,
	}

	// Create confirmation response
	response, err := accounts.NewConfirmationResponse(testConfirmation)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, fmt.Sprintf("/v1/confirmations/%d", testConfirmation.ID), selfLink.Href)
	}
}

func TestNewInvitationResponse(t *testing.T) {
	testInvitedUser := &accounts.User{
		OauthUser: new(oauth.User),
	}
	testInvitedByUser := &accounts.User{
		OauthUser: new(oauth.User),
	}
	testInvitation := &accounts.Invitation{
		InvitedUser:   testInvitedUser,
		InvitedByUser: testInvitedByUser,
	}

	// Create invitation response
	response, err := accounts.NewInvitationResponse(testInvitation)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, fmt.Sprintf("/v1/invitations/%d", testInvitation.ID), selfLink.Href)
	}
}

func TestNewPasswordResetResponse(t *testing.T) {
	testUser := &accounts.User{
		OauthUser: new(oauth.User),
	}
	testPasswordReset := &accounts.PasswordReset{
		User: testUser,
	}

	// Create password reset response
	response, err := accounts.NewPasswordResetResponse(testPasswordReset)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, fmt.Sprintf("/v1/password-resets/%d", testPasswordReset.ID), selfLink.Href)
	}
}

func TestNewListUsersResponse(t *testing.T) {
	testUsers := []*accounts.User{
		&accounts.User{
			OauthUser: new(oauth.User),
		},
		&accounts.User{
			OauthUser: new(oauth.User),
		},
	}

	// Create list response
	response, err := accounts.NewListUsersResponse(
		10,                 // count
		2,                  // page
		"/v1/users?page=2", // self
		"/v1/users?page=1", // first
		"/v1/users?page=5", // last
		"/v1/users?page=1", // previous
		"/v1/users?page=3", // next
		testUsers,
	)

	// Error should be nil
	assert.Nil(t, err)

	// Test self link
	selfLink, err := response.GetLink("self")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/users?page=2", selfLink.Href)
	}

	// Test first link
	firstLink, err := response.GetLink("first")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/users?page=1", firstLink.Href)
	}

	// Test last link
	lastLink, err := response.GetLink("last")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/users?page=5", lastLink.Href)
	}

	// Test previous link
	previousLink, err := response.GetLink("prev")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/users?page=1", previousLink.Href)
	}

	// Test next link
	nextLink, err := response.GetLink("next")
	if assert.Nil(t, err) {
		assert.Equal(t, "/v1/users?page=3", nextLink.Href)
	}

	// Test embedded users
	embeddedUsers, err := response.GetEmbedded("users")
	if assert.Nil(t, err) {
		reflectedValue := reflect.ValueOf(embeddedUsers)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(accounts.UserResponse)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}

	// Test the rest
	assert.Equal(t, uint(10), response.Count)
	assert.Equal(t, uint(2), response.Page)
}
