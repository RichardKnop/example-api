package accounts

import (
	"github.com/RichardKnop/recall/oauth"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindAccountByOauthClientID() {
	var (
		account *Account
		err     error
	)

	// Let's try to find an account by a bogus ID
	account, err = suite.service.FindAccountByOauthClientID(12345)

	// Account should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccountNotFound, err)
	}

	// Now let's pass a valid ID
	account, err = suite.service.FindAccountByOauthClientID(suite.accounts[0].OauthClient.ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account should be returned
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account 1", account.Name)
		assert.Equal(suite.T(), "test_client_1", account.OauthClient.Key)
	}
}

func (suite *AccountsTestSuite) TestFindAccountByID() {
	var (
		account *Account
		err     error
	)

	// Let's try to find an account by a bogus ID
	account, err = suite.service.FindAccountByID(12345)

	// Account should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccountNotFound, err)
	}

	// Now let's pass a valid ID
	account, err = suite.service.FindAccountByID(suite.accounts[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account should be returned with preloaded data
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account 1", account.Name)
		assert.Equal(suite.T(), "test_client_1", account.OauthClient.Key)
	}
}

func (suite *AccountsTestSuite) TestFindAccountByName() {
	var (
		account *Account
		err     error
	)

	// Let's try to find an account by a bogus name
	account, err = suite.service.FindAccountByName("bogus")

	// Account should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccountNotFound, err)
	}

	// Now let's pass a valid name
	account, err = suite.service.FindAccountByName(suite.accounts[0].Name)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account should be returned with preloaded data
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account 1", account.Name)
		assert.Equal(suite.T(), "test_client_1", account.OauthClient.Key)
	}
}

func (suite *AccountsTestSuite) TestCreateAccount() {
	var (
		account *Account
		err     error
	)

	// We try to insert an account with a non unique oauth client
	account, err = suite.service.CreateAccount(
		"New Account",             // name
		"",                        // description
		"test_client_2",           // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Account object should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrClientIDTaken, err)
	}

	// We try to insert an account with a non unique name
	account, err = suite.service.CreateAccount(
		"Test Account 2",          // name
		"",                        // description
		"new_client",              // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Account object should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAccountNameTaken, err)
	}

	// We try to insert a unique account
	account, err = suite.service.CreateAccount(
		"New Account",             // name
		"",                        // description
		"new_client",              // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account object should be returned
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "New Account", account.Name)
		assert.Equal(suite.T(), "new_client", account.OauthClient.Key)
	}
}
