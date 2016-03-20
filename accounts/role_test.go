package accounts

import (
	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindRoleByID() {
	var (
		role *Role
		err  error
	)

	// Let's try to find a role by a bogus ID
	role, err = suite.service.findRoleByID("bogus")

	// Role should be nil
	assert.Nil(suite.T(), role)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrRoleNotFound, err)
	}

	// Now let's pass a valid ID
	role, err = suite.service.findRoleByID(roles.User)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct role should be returned
	if assert.NotNil(suite.T(), role) {
		assert.Equal(suite.T(), roles.User, role.ID)
	}
}
