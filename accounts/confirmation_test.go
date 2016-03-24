package accounts

import (
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindConfirmationByReference() {
	var (
		confirmation *Confirmation
		err          error
	)

	testConfirmation := NewConfirmation(suite.users[1])
	err = suite.db.Create(testConfirmation).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Let's try to find an confirmation by a bogus reference
	confirmation, err = suite.service.FindConfirmationByReference("bogus")

	// Confirmation should be nil
	assert.Nil(suite.T(), confirmation)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrConfirmationNotFound, err)
	}

	// Now let's pass a valid reference
	confirmation, err = suite.service.FindConfirmationByReference(testConfirmation.Reference)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct confirmation should be returned with preloaded data
	if assert.NotNil(suite.T(), confirmation) {
		assert.Equal(suite.T(), suite.users[1].ID, confirmation.User.ID)
		assert.False(suite.T(), confirmation.EmailSent)
		assert.False(suite.T(), confirmation.EmailSentAt.Valid)
	}
}
