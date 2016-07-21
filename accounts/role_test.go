package accounts_test

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/accounts/roles"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestFindRoleByID() {
	var (
		role *accounts.Role
		err  error
	)

	// Let's try to find a role by a bogus ID
	role, err = suite.service.FindRoleByID("bogus")

	// Role should be nil
	assert.Nil(suite.T(), role)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), accounts.ErrRoleNotFound, err)
	}

	// Now let's pass a valid ID
	role, err = suite.service.FindRoleByID(roles.User)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct role should be returned
	if assert.NotNil(suite.T(), role) {
		assert.Equal(suite.T(), roles.User, role.ID)
	}
}
