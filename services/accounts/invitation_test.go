package accounts_test

import (
	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindInvitationByReference() {
	var (
		testExpiredInvitation, testInvitation, invitation *models.Invitation
		err                                               error
	)

	// Insert test invitations

	testExpiredInvitation, err = models.NewInvitation(
		suite.users[0],
		suite.users[1],
		-10, // expires in
	)
	assert.NoError(suite.T(), err, "Failed to create a new invitation object")
	err = suite.db.Create(testExpiredInvitation).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	testInvitation, err = models.NewInvitation(
		suite.users[0],
		suite.users[1],
		suite.cnf.AppSpecific.InvitationLifetime,
	)
	assert.NoError(suite.T(), err, "Failed to create a new invitation object")
	err = suite.db.Create(testInvitation).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Let's try to find an expired invitation by a valid reference
	invitation, err = suite.service.FindInvitationByReference(testExpiredInvitation.Reference)

	// Invitation should be nil
	assert.Nil(suite.T(), invitation)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrInvitationNotFound, err)
	}

	// Let's try to find an invitation by a bogus reference
	invitation, err = suite.service.FindInvitationByReference("bogus")

	// Invitation should be nil
	assert.Nil(suite.T(), invitation)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrInvitationNotFound, err)
	}

	// Now let's pass a valid reference
	invitation, err = suite.service.FindInvitationByReference(testInvitation.Reference)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct invitation should be returned with preloaded data
	if assert.NotNil(suite.T(), invitation) {
		assert.Equal(suite.T(), suite.users[0].ID, invitation.InvitedUser.ID)
		assert.Equal(suite.T(), suite.users[1].ID, invitation.InvitedByUser.ID)
		assert.False(suite.T(), invitation.EmailSent)
		assert.Nil(suite.T(), invitation.EmailSentAt)
	}
}
