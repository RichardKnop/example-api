package accounts

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindPasswordResetByReference() {
	var (
		passwordReset *PasswordReset
		validFor      = time.Duration(suite.service.cnf.Recall.PasswordResetLifetime) * time.Second
		err           error
	)

	// Insert a test password reset
	testPasswordReset := NewPasswordReset(suite.users[1])
	err = suite.db.Create(testPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test password reset failed")
	err = suite.db.Model(testPasswordReset).UpdateColumn(
		"created_at",
		time.Now().Add(-validFor).Add(time.Second),
	).Error
	assert.NoError(suite.T(), err, "Updating test password reset failed")

	// Insert a test expired password reset
	testExpiredPasswordReset := NewPasswordReset(suite.users[1])
	err = suite.db.Create(testExpiredPasswordReset).Error
	assert.NoError(suite.T(), err, "Inserting test expired password reset failed")
	err = suite.db.Model(testExpiredPasswordReset).UpdateColumn(
		"created_at",
		time.Now().Add(-validFor).Add(-time.Second),
	).Error
	assert.NoError(suite.T(), err, "Updating test expired password reset failed")

	// Let's try to find a password reset by a bogus reference
	passwordReset, err = suite.service.FindPasswordResetByReference("bogus")

	// Password reset should be nil
	assert.Nil(suite.T(), passwordReset)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrPasswordResetNotFound, err)
	}

	// Now let's pass a valid reference of an expired password reset
	passwordReset, err = suite.service.FindPasswordResetByReference(testExpiredPasswordReset.Reference)

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
