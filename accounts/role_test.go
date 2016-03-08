package accounts

import (
	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindRoleByName() {
	var (
		role *Role
		err  error
	)

	role, err = suite.service.findRoleByName("bogus")

	// Role should be nil
	assert.Nil(suite.T(), role)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrRoleNotFound, err)
	}

	role, err = suite.service.findRoleByName(roles.User)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct role should be returned
	if assert.NotNil(suite.T(), role) {
		assert.Equal(suite.T(), roles.User, role.Name)
	}
}
