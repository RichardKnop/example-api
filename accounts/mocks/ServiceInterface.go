package mocks

import "github.com/RichardKnop/example-api/accounts"
import "github.com/stretchr/testify/mock"

import "net/http"
import "github.com/RichardKnop/example-api/config"
import "github.com/RichardKnop/example-api/models"
import "github.com/RichardKnop/example-api/oauth"
import "github.com/RichardKnop/example-api/util/routes"
import "github.com/gorilla/mux"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) GetConfig() *config.Config {
	ret := _m.Called()

	var r0 *config.Config
	if rf, ok := ret.Get(0).(func() *config.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Config)
		}
	}

	return r0
}
func (_m *ServiceInterface) GetOauthService() oauth.ServiceInterface {
	ret := _m.Called()

	var r0 oauth.ServiceInterface
	if rf, ok := ret.Get(0).(func() oauth.ServiceInterface); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(oauth.ServiceInterface)
	}

	return r0
}
func (_m *ServiceInterface) GetRoutes() []routes.Route {
	ret := _m.Called()

	var r0 []routes.Route
	if rf, ok := ret.Get(0).(func() []routes.Route); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]routes.Route)
		}
	}

	return r0
}
func (_m *ServiceInterface) RegisterRoutes(router *mux.Router, prefix string) {
	_m.Called(router, prefix)
}
func (_m *ServiceInterface) GetUserCredentialsFromToken(token string) (*models.User, error) {
	ret := _m.Called(token)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetClientCredentialsFromBaseAuth(r *http.Request) (*models.Account, error) {
	ret := _m.Called(r)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(*http.Request) *models.Account); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetClientCredentialsFromToken(token string) (*models.Account, error) {
	ret := _m.Called(token)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(string) *models.Account); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) GetMixedCredentialsFromToken(token string) (*models.Account, *models.User, error) {
	ret := _m.Called(token)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(string) *models.Account); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 *models.User
	if rf, ok := ret.Get(1).(func(string) *models.User); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*models.User)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(token)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ServiceInterface) FindAccountByOauthClientID(oauthClientID uint) (*models.Account, error) {
	ret := _m.Called(oauthClientID)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(uint) *models.Account); ok {
		r0 = rf(oauthClientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(oauthClientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindAccountByID(accountID uint) (*models.Account, error) {
	ret := _m.Called(accountID)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(uint) *models.Account); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindAccountByName(name string) (*models.Account, error) {
	ret := _m.Called(name)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(string) *models.Account); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateAccount(name string, description string, key string, secret string, redirectURI string) (*models.Account, error) {
	ret := _m.Called(name, description, key, secret, redirectURI)

	var r0 *models.Account
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) *models.Account); ok {
		r0 = rf(name, description, key, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string) error); ok {
		r1 = rf(name, description, key, secret, redirectURI)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByOauthUserID(oauthUserID uint) (*models.User, error) {
	ret := _m.Called(oauthUserID)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(uint) *models.User); ok {
		r0 = rf(oauthUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(oauthUserID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByEmail(email string) (*models.User, error) {
	ret := _m.Called(email)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByID(userID uint) (*models.User, error) {
	ret := _m.Called(userID)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(uint) *models.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindUserByFacebookID(facebookID string) (*models.User, error) {
	ret := _m.Called(facebookID)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(string) *models.User); ok {
		r0 = rf(facebookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(facebookID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateUser(account *models.Account, userRequest *accounts.UserRequest) (*models.User, error) {
	ret := _m.Called(account, userRequest)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*models.Account, *accounts.UserRequest) *models.User); ok {
		r0 = rf(account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Account, *accounts.UserRequest) error); ok {
		r1 = rf(account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) UpdateUser(user *models.User, userRequest *accounts.UserRequest) error {
	ret := _m.Called(user, userRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User, *accounts.UserRequest) error); ok {
		r0 = rf(user, userRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) PaginatedUsersCount() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindPaginatedUsers(offset int, limit int, sorts map[string]string) ([]*models.User, error) {
	ret := _m.Called(offset, limit, sorts)

	var r0 []*models.User
	if rf, ok := ret.Get(0).(func(int, int, map[string]string) []*models.User); ok {
		r0 = rf(offset, limit, sorts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, map[string]string) error); ok {
		r1 = rf(offset, limit, sorts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindConfirmationByReference(reference string) (*models.Confirmation, error) {
	ret := _m.Called(reference)

	var r0 *models.Confirmation
	if rf, ok := ret.Get(0).(func(string) *models.Confirmation); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Confirmation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(reference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) ConfirmUser(confirmation *models.Confirmation) error {
	ret := _m.Called(confirmation)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Confirmation) error); ok {
		r0 = rf(confirmation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) FindPasswordResetByReference(reference string) (*models.PasswordReset, error) {
	ret := _m.Called(reference)

	var r0 *models.PasswordReset
	if rf, ok := ret.Get(0).(func(string) *models.PasswordReset); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PasswordReset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(reference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) ResetPassword(passwordReset *models.PasswordReset, password string) error {
	ret := _m.Called(passwordReset, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.PasswordReset, string) error); ok {
		r0 = rf(passwordReset, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) GetOrCreateFacebookUser(account *models.Account, facebookID string, userRequest *accounts.UserRequest) (*models.User, error) {
	ret := _m.Called(account, facebookID, userRequest)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*models.Account, string, *accounts.UserRequest) *models.User); ok {
		r0 = rf(account, facebookID, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Account, string, *accounts.UserRequest) error); ok {
		r1 = rf(account, facebookID, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateSuperuser(account *models.Account, email string, password string) (*models.User, error) {
	ret := _m.Called(account, email, password)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*models.Account, string, string) *models.User); ok {
		r0 = rf(account, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Account, string, string) error); ok {
		r1 = rf(account, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindInvitationByReference(reference string) (*models.Invitation, error) {
	ret := _m.Called(reference)

	var r0 *models.Invitation
	if rf, ok := ret.Get(0).(func(string) *models.Invitation); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invitation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(reference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) InviteUser(invitedByUser *models.User, invitationRequest *accounts.InvitationRequest) (*models.Invitation, error) {
	ret := _m.Called(invitedByUser, invitationRequest)

	var r0 *models.Invitation
	if rf, ok := ret.Get(0).(func(*models.User, *accounts.InvitationRequest) *models.Invitation); ok {
		r0 = rf(invitedByUser, invitationRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Invitation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User, *accounts.InvitationRequest) error); ok {
		r1 = rf(invitedByUser, invitationRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) ConfirmInvitation(invitation *models.Invitation, password string) error {
	ret := _m.Called(invitation, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Invitation, string) error); ok {
		r0 = rf(invitation, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) GetUserFromQueryString(r *http.Request) (*models.User, error) {
	ret := _m.Called(r)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(*http.Request) *models.User); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
