package mocks

import "github.com/RichardKnop/example-api/accounts"
import "github.com/stretchr/testify/mock"

import "net/http"
import "github.com/RichardKnop/example-api/config"
import "github.com/RichardKnop/example-api/oauth"
import "github.com/jinzhu/gorm"

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
func (_m *ServiceInterface) FindAccountByOauthClientID(oauthClientID uint) (*accounts.Account, error) {
	ret := _m.Called(oauthClientID)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(uint) *accounts.Account); ok {
		r0 = rf(oauthClientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
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
func (_m *ServiceInterface) FindAccountByID(accountID uint) (*accounts.Account, error) {
	ret := _m.Called(accountID)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(uint) *accounts.Account); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
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
func (_m *ServiceInterface) FindAccountByName(name string) (*accounts.Account, error) {
	ret := _m.Called(name)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(string) *accounts.Account); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
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
func (_m *ServiceInterface) CreateAccount(name string, description string, key string, secret string, redirectURI string) (*accounts.Account, error) {
	ret := _m.Called(name, description, key, secret, redirectURI)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) *accounts.Account); ok {
		r0 = rf(name, description, key, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
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
func (_m *ServiceInterface) FindUserByOauthUserID(oauthUserID uint) (*accounts.User, error) {
	ret := _m.Called(oauthUserID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(uint) *accounts.User); ok {
		r0 = rf(oauthUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
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
func (_m *ServiceInterface) FindUserByEmail(email string) (*accounts.User, error) {
	ret := _m.Called(email)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(string) *accounts.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
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
func (_m *ServiceInterface) FindUserByID(userID uint) (*accounts.User, error) {
	ret := _m.Called(userID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(uint) *accounts.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
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
func (_m *ServiceInterface) FindUserByFacebookID(facebookID string) (*accounts.User, error) {
	ret := _m.Called(facebookID)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(string) *accounts.User); ok {
		r0 = rf(facebookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
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
func (_m *ServiceInterface) CreateUser(account *accounts.Account, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(account, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, *accounts.UserRequest) error); ok {
		r1 = rf(account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateUserTx(tx *gorm.DB, account *accounts.Account, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(tx, account, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*gorm.DB, *accounts.Account, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(tx, account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *accounts.Account, *accounts.UserRequest) error); ok {
		r1 = rf(tx, account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) UpdateUser(user *accounts.User, userRequest *accounts.UserRequest) error {
	ret := _m.Called(user, userRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.User, *accounts.UserRequest) error); ok {
		r0 = rf(user, userRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) FindConfirmationByReference(reference string) (*accounts.Confirmation, error) {
	ret := _m.Called(reference)

	var r0 *accounts.Confirmation
	if rf, ok := ret.Get(0).(func(string) *accounts.Confirmation); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Confirmation)
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
func (_m *ServiceInterface) ConfirmUser(user *accounts.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) FindPasswordResetByReference(reference string) (*accounts.PasswordReset, error) {
	ret := _m.Called(reference)

	var r0 *accounts.PasswordReset
	if rf, ok := ret.Get(0).(func(string) *accounts.PasswordReset); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.PasswordReset)
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
func (_m *ServiceInterface) ResetPassword(passwordReset *accounts.PasswordReset, password string) error {
	ret := _m.Called(passwordReset, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.PasswordReset, string) error); ok {
		r0 = rf(passwordReset, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) GetOrCreateFacebookUser(account *accounts.Account, facebookID string, userRequest *accounts.UserRequest) (*accounts.User, error) {
	ret := _m.Called(account, facebookID, userRequest)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, string, *accounts.UserRequest) *accounts.User); ok {
		r0 = rf(account, facebookID, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, string, *accounts.UserRequest) error); ok {
		r1 = rf(account, facebookID, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) CreateSuperuser(account *accounts.Account, email string, password string) (*accounts.User, error) {
	ret := _m.Called(account, email, password)

	var r0 *accounts.User
	if rf, ok := ret.Get(0).(func(*accounts.Account, string, string) *accounts.User); ok {
		r0 = rf(account, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Account, string, string) error); ok {
		r1 = rf(account, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindInvitationByID(invitationID uint) (*accounts.Invitation, error) {
	ret := _m.Called(invitationID)

	var r0 *accounts.Invitation
	if rf, ok := ret.Get(0).(func(uint) *accounts.Invitation); ok {
		r0 = rf(invitationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Invitation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(invitationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) FindInvitationByReference(reference string) (*accounts.Invitation, error) {
	ret := _m.Called(reference)

	var r0 *accounts.Invitation
	if rf, ok := ret.Get(0).(func(string) *accounts.Invitation); ok {
		r0 = rf(reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Invitation)
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
func (_m *ServiceInterface) InviteUser(invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	ret := _m.Called(invitedByUser, invitationRequest)

	var r0 *accounts.Invitation
	if rf, ok := ret.Get(0).(func(*accounts.User, *accounts.InvitationRequest) *accounts.Invitation); ok {
		r0 = rf(invitedByUser, invitationRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Invitation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.User, *accounts.InvitationRequest) error); ok {
		r1 = rf(invitedByUser, invitationRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) InviteUserTx(tx *gorm.DB, invitedByUser *accounts.User, invitationRequest *accounts.InvitationRequest) (*accounts.Invitation, error) {
	ret := _m.Called(tx, invitedByUser, invitationRequest)

	var r0 *accounts.Invitation
	if rf, ok := ret.Get(0).(func(*gorm.DB, *accounts.User, *accounts.InvitationRequest) *accounts.Invitation); ok {
		r0 = rf(tx, invitedByUser, invitationRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Invitation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gorm.DB, *accounts.User, *accounts.InvitationRequest) error); ok {
		r1 = rf(tx, invitedByUser, invitationRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ServiceInterface) ConfirmInvitation(invitation *accounts.Invitation, password string) error {
	ret := _m.Called(invitation, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(*accounts.Invitation, string) error); ok {
		r0 = rf(invitation, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ServiceInterface) GetUserCredentials(token string) (*accounts.Account, *accounts.User, error) {
	ret := _m.Called(token)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(string) *accounts.Account); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 *accounts.User
	if rf, ok := ret.Get(1).(func(string) *accounts.User); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*accounts.User)
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
func (_m *ServiceInterface) GetClientCredentials(r *http.Request) (*accounts.Account, *accounts.User, error) {
	ret := _m.Called(r)

	var r0 *accounts.Account
	if rf, ok := ret.Get(0).(func(*http.Request) *accounts.Account); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*accounts.Account)
		}
	}

	var r1 *accounts.User
	if rf, ok := ret.Get(1).(func(*http.Request) *accounts.User); ok {
		r1 = rf(r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*accounts.User)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*http.Request) error); ok {
		r2 = rf(r)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
func (_m *ServiceInterface) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) GetMyUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) InviteUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) CreatePasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) ConfirmInvitationHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
func (_m *ServiceInterface) ConfirmPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
