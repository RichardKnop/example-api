package mocks

import "github.com/RichardKnop/example-api/email"
import "github.com/stretchr/testify/mock"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) Send(e *email.Email) error {
	ret := _m.Called(e)

	var r0 error
	if rf, ok := ret.Get(0).(func(*email.Email) error); ok {
		r0 = rf(e)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
