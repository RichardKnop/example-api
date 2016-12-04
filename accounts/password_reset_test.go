package accounts_test

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindPasswordResetByReference() {
	var (
		testExpiredPasswordReset, testPasswordReset, passwordReset *models.PasswordReset
		err                                                        error
	)

	// Insert test password resets

	testExpiredPasswordReset, err = models.NewPasswordReset(
		suite.users[1],
		-10, // expires in
	)
	assert.NoError(suite.T(), err, "Failed to create a new password reset object")
	err = suite.db.Create(testExpiredPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test expired password reset failed")
	testExpiredPasswordReset.User = suite.users[1]

	testPasswordReset, err = models.NewPasswordReset(
		suite.users[1],
		suite.cnf.AppSpecific.PasswordResetLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new password reset object")
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test password reset failed")
	testPasswordReset.User = suite.users[1]

	// Let's try to find an expired password reset by a valid reference
	passwordReset, err = suite.service.FindPasswordResetByReference(testExpiredPasswordReset.Reference)

	// Password reset should be nil
	assert.Nil(suite.T(), passwordReset)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrPasswordResetNotFound, err)
	}

	// Let's try to find a password reset by a bogus reference
	passwordReset, err = suite.service.FindPasswordResetByReference("bogus")

	// Password reset should be nil
	assert.Nil(suite.T(), passwordReset)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrPasswordResetNotFound, err)
	}

	// Now let's pass a valid reference of an expired password reset
	passwordReset, err = suite.service.FindPasswordResetByReference(testExpiredPasswordReset.Reference)

	// Password reset should be nil
	assert.Nil(suite.T(), passwordReset)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrPasswordResetNotFound, err)
	}

	// Now let's pass a valid reference
	passwordReset, err = suite.service.FindPasswordResetByReference(testPasswordReset.Reference)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct password reset should be returned with preloaded data
	if assert.NotNil(suite.T(), passwordReset) {
		assert.Equal(suite.T(), suite.users[1].ID, passwordReset.User.ID)
		assert.False(suite.T(), passwordReset.EmailSent)
		assert.Nil(suite.T(), passwordReset.EmailSentAt)
	}
}

func (suite *AccountsTestSuite) TestResetPassword() {
	// Insert a test password reset
	testPasswordReset, err := models.NewPasswordReset(
		suite.users[1],
		suite.cnf.AppSpecific.PasswordResetLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new password reset object")
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")
	testPasswordReset.User = suite.users[1]

	// Error should be nil
	assert.Nil(
		suite.T(),
		suite.service.ResetPassword(testPasswordReset, "newpassword"),
	)

	// The password reset object should have been deleted
	assert.True(
		suite.T(),
		suite.db.Find(new(models.PasswordReset), testPasswordReset.ID).RecordNotFound(),
	)
}
