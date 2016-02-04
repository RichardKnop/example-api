package accounts

import (
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
		assert.Equal(suite.T(), errAccountNotFound, err)
	}

	// Now let's pass a valid ID
	account, err = suite.service.FindAccountByOauthClientID(suite.accounts[0].OauthClient.ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account should be returned
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account", account.Name)
		assert.Equal(suite.T(), "test_client", account.OauthClient.Key)
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
		assert.Equal(suite.T(), errAccountNotFound, err)
	}

	// Now let's pass a valid ID
	account, err = suite.service.FindAccountByID(suite.accounts[0].ID)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account should be returned with preloaded data
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account", account.Name)
		assert.Equal(suite.T(), "test_client", account.OauthClient.Key)
	}
}

func (suite *AccountsTestSuite) TestCreateAccount() {
	var (
		account *Account
		err     error
	)

	// We try to insert an account with a non unique oauth client
	account, err = suite.service.CreateAccount(
		"Test Account 2",          // name
		"",                        // description
		"test_client",             // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Account object should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_clients.key", err.Error())
	}

	// We try to insert a non unique account
	account, err = suite.service.CreateAccount(
		"Test Account",            // name
		"",                        // description
		"test_client2",            // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Account object should be nil
	assert.Nil(suite.T(), account)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: account_accounts.name", err.Error())
	}

	// We try to insert a unique account
	account, err = suite.service.CreateAccount(
		"Test Account 2",          // name
		"",                        // description
		"test_client2",            // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct account object should be returned
	if assert.NotNil(suite.T(), account) {
		assert.Equal(suite.T(), "Test Account 2", account.Name)
		assert.Equal(suite.T(), "test_client2", account.OauthClient.Key)
	}
}
