package mocks

import "github.com/RichardKnop/example-api/accounts"
import "github.com/stretchr/testify/mock"

import "github.com/RichardKnop/example-api/email"

type EmailFactoryInterface struct {
	mock.Mock
}

func (_m *EmailFactoryInterface) NewConfirmationEmail(o *accounts.Confirmation) (*email.Message, error) {
	ret := _m.Called(o)

	var r0 *email.Message
	if rf, ok := ret.Get(0).(func(*accounts.Confirmation) *email.Message); ok {
		r0 = rf(o)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Confirmation) error); ok {
		r1 = rf(o)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *EmailFactoryInterface) NewInvitationEmail(o *accounts.Invitation) (*email.Message, error) {
	ret := _m.Called(o)

	var r0 *email.Message
	if rf, ok := ret.Get(0).(func(*accounts.Invitation) *email.Message); ok {
		r0 = rf(o)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.Invitation) error); ok {
		r1 = rf(o)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *EmailFactoryInterface) NewPasswordResetEmail(o *accounts.PasswordReset) (*email.Message, error) {
	ret := _m.Called(o)

	var r0 *email.Message
	if rf, ok := ret.Get(0).(func(*accounts.PasswordReset) *email.Message); ok {
		r0 = rf(o)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*email.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*accounts.PasswordReset) error); ok {
		r1 = rf(o)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
