package accounts_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/example-api/accounts"
)

func (suite *AccountsTestSuite) TestFindConfirmationByReference() {
	var (
		testExpiredConfirmation, testConfirmation, confirmation *accounts.Confirmation
		err                                                     error
	)

	// Insert test confirmations

	testExpiredConfirmation, err = accounts.NewConfirmation(
		suite.users[1],
		-10, // expires in
	)
	assert.NoError(suite.T(), err, "Failed to create a new confirmation object")
	err = suite.db.Create(testExpiredConfirmation).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	testConfirmation, err = accounts.NewConfirmation(
		suite.users[1],
		suite.cnf.AppSpecific.ConfirmationLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new confirmation object")
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Let's try to find an expired confirmation by a valid reference
	confirmation, err = suite.service.FindConfirmationByReference(testExpiredConfirmation.Reference)

	// Confirmation should be nil
	assert.Nil(suite.T(), confirmation)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrConfirmationNotFound, err)
	}

	// Let's try to find a confirmation by a bogus reference
	confirmation, err = suite.service.FindConfirmationByReference("bogus")

	// Confirmation should be nil
	assert.Nil(suite.T(), confirmation)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrConfirmationNotFound, err)
	}

	// Now let's pass a valid reference
	confirmation, err = suite.service.FindConfirmationByReference(testConfirmation.Reference)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct confirmation should be returned with preloaded data
	if assert.NotNil(suite.T(), confirmation) {
		assert.Equal(suite.T(), suite.users[1].ID, confirmation.User.ID)
		assert.False(suite.T(), confirmation.EmailSent)
		assert.Nil(suite.T(), confirmation.EmailSentAt)
	}
}
