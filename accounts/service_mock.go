package accounts

import (
	"net/http"

	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/oauth"
	"github.com/stretchr/testify/mock"
)

// ServiceMock is a mocked object implementing ServiceInterface
type ServiceMock struct {
	mock.Mock
}

// GetConfig ...
func (_m *ServiceMock) GetConfig() *config.Config {
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

// GetOauthService ...
func (_m *ServiceMock) GetOauthService() oauth.ServiceInterface {
	ret := _m.Called()

	var r0 oauth.ServiceInterface
	if rf, ok := ret.Get(0).(func() oauth.ServiceInterface); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(oauth.ServiceInterface)
	}

	return r0
}

// FindAccountByOauthClientID ...
func (_m *ServiceMock) FindAccountByOauthClientID(oauthClientID uint) (*Account, error) {
	ret := _m.Called(oauthClientID)

	var r0 *Account
	if rf, ok := ret.Get(0).(func(uint) *Account); ok {
		r0 = rf(oauthClientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Account)
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

// FindAccountByID ...
func (_m *ServiceMock) FindAccountByID(accountID uint) (*Account, error) {
	ret := _m.Called(accountID)

	var r0 *Account
	if rf, ok := ret.Get(0).(func(uint) *Account); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Account)
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

// CreateAccount ...
func (_m *ServiceMock) CreateAccount(name string, description string, key string, secret string, redirectURI string) (*Account, error) {
	ret := _m.Called(name, description, key, secret, redirectURI)

	var r0 *Account
	if rf, ok := ret.Get(0).(func(string, string, string, string, string) *Account); ok {
		r0 = rf(name, description, key, secret, redirectURI)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Account)
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

// FindUserByOauthUserID ...
func (_m *ServiceMock) FindUserByOauthUserID(oauthUserID uint) (*User, error) {
	ret := _m.Called(oauthUserID)

	var r0 *User
	if rf, ok := ret.Get(0).(func(uint) *User); ok {
		r0 = rf(oauthUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
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

// FindUserByID ...
func (_m *ServiceMock) FindUserByID(userID uint) (*User, error) {
	ret := _m.Called(userID)

	var r0 *User
	if rf, ok := ret.Get(0).(func(uint) *User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
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

// FindUserByFacebookID ...
func (_m *ServiceMock) FindUserByFacebookID(facebookID string) (*User, error) {
	ret := _m.Called(facebookID)

	var r0 *User
	if rf, ok := ret.Get(0).(func(string) *User); ok {
		r0 = rf(facebookID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
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

// CreateUser ...
func (_m *ServiceMock) CreateUser(account *Account, userRequest *UserRequest) (*User, error) {
	ret := _m.Called(account, userRequest)

	var r0 *User
	if rf, ok := ret.Get(0).(func(*Account, *UserRequest) *User); ok {
		r0 = rf(account, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Account, *UserRequest) error); ok {
		r1 = rf(account, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser ...
func (_m *ServiceMock) UpdateUser(user *User, userRequest *UserRequest) error {
	ret := _m.Called(user, userRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*User, *UserRequest) error); ok {
		r0 = rf(user, userRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateFacebookUser ...
func (_m *ServiceMock) CreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error) {
	ret := _m.Called(account, facebookID, userRequest)

	var r0 *User
	if rf, ok := ret.Get(0).(func(*Account, string, *UserRequest) *User); ok {
		r0 = rf(account, facebookID, userRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Account, string, *UserRequest) error); ok {
		r1 = rf(account, facebookID, userRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSuperuser ..
func (_m *ServiceMock) CreateSuperuser(account *Account, email string, password string) (*User, error) {
	ret := _m.Called(account, email, password)

	var r0 *User
	if rf, ok := ret.Get(0).(func(*Account, string, string) *User); ok {
		r0 = rf(account, email, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Account, string, string) error); ok {
		r1 = rf(account, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *ServiceMock) createUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) getMyUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

func (_m *ServiceMock) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}
