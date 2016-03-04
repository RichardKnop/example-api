package facebook

import (
	fb "github.com/huandu/facebook"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// AdapterMock is a mocked object implementing AdapterInterface
type AdapterMock struct {
	mock.Mock
}

// AuthCodeURL ...
func (_m *AdapterMock) AuthCodeURL(state string) string {
	ret := _m.Called(state)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Exchange ...
func (_m *AdapterMock) Exchange(code string) (*oauth2.Token, error) {
	ret := _m.Called(code)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(string) *oauth2.Token); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMe ...
func (_m *AdapterMock) GetMe(accessToken string) (fb.Result, error) {
	ret := _m.Called(accessToken)

	var r0 fb.Result
	if rf, ok := ret.Get(0).(func(string) fb.Result); ok {
		r0 = rf(accessToken)
	} else {
		r0 = ret.Get(0).(fb.Result)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
