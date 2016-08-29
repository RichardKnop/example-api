package email

// ServiceInterface defines exported methods
type ServiceInterface interface {
	Send(m *Message) error
}
