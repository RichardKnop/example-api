package email

// ServiceInterface defines exported methods
type ServiceInterface interface {
	Send(email *Email) error
}
