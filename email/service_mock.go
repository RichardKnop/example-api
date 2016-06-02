package email

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// ServiceMock is a mocked object implementing ServiceInterface
type ServiceMock struct {
	mock.Mock
}

// Send just records the activity, and returns what the Mock object tells it to
func (m *ServiceMock) Send(email *Email) error {
	time.Sleep(10 * time.Millisecond)
	args := m.Called(email)
	return args.Error(0)
}
