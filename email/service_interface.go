package email

// ServiceInterface defines exported methods
type ServiceInterface interface {
	Send(e *Email) error
}
