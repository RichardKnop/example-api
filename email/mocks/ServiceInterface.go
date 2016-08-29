package mocks

import "github.com/RichardKnop/example-api/email"
import "github.com/stretchr/testify/mock"

type ServiceInterface struct {
	mock.Mock
}

func (_m *ServiceInterface) Send(m *email.Message) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*email.Message) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
