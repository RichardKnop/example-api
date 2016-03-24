package accounts

import (
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindPasswordResetByReference() {
	var (
		passwordReset *PasswordReset
		err           error
	)

	// Insert a test password reset
	testPasswordReset := NewPasswordReset(suite.users[1])
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Let's try to find a password reset by a bogus reference
	passwordReset, err = suite.service.FindPasswordResetByReference("bogus")

	// Password reset should be nil
	assert.Nil(suite.T(), passwordReset)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrPasswordResetNotFound, err)
	}

	// Now let's pass a valid reference
	passwordReset, err = suite.service.FindPasswordResetByReference(testPasswordReset.Reference)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct passwor dreset should be returned with preloaded data
	if assert.NotNil(suite.T(), passwordReset) {
		assert.Equal(suite.T(), suite.users[1].ID, passwordReset.User.ID)
		assert.False(suite.T(), passwordReset.EmailSent)
		assert.False(suite.T(), passwordReset.EmailSentAt.Valid)
	}
}

func (suite *AccountsTestSuite) TestResetPassword() {
	// Insert a test password reset
	passwordReset := NewPasswordReset(suite.users[1])
	err := suite.db.Create(passwordReset).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Error should be nil
	assert.Nil(
		suite.T(),
		suite.service.ResetPassword(passwordReset, "newpassword"),
	)

	// The password reset object should have been deleted
	assert.True(
		suite.T(),
		suite.db.Find(new(PasswordReset), passwordReset.ID).RecordNotFound(),
	)
}
