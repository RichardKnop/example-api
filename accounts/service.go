package accounts

import (
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf           *config.Config
	db            *gorm.DB
	oauthService  oauth.ServiceInterface
	emailService  email.ServiceInterface
	emailFactory  EmailFactoryInterface
	notifications chan bool
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, oauthService oauth.ServiceInterface, emailService email.ServiceInterface, emailFactory EmailFactoryInterface) *Service {
	if emailFactory == nil {
		emailFactory = NewEmailFactory(cnf)
	}
	return &Service{
		cnf:          cnf,
		db:           db,
		oauthService: oauthService,
		emailService: emailService,
		emailFactory: emailFactory,
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// GetOauthService returns oauth.Service instance
func (s *Service) GetOauthService() oauth.ServiceInterface {
	return s.oauthService
}

// WaitForNotifications informs the service to expect a number (size) of Notify() requests.
// This is useful for testing to ensure all async tasks have finished
func (s *Service) WaitForNotifications(size int) {
	s.notifications = make(chan bool, size)
}

// Notify increments the number of notifications received for the service.
// This is useful for testing to ensure all async tasks have finished
func (s *Service) Notify() {
	if s.notifications != nil {
		s.notifications <- true
	}
}

// DontWaitForNotifications informs the service not to wait for any notifications
// This is useful for testing to ensure all async tasks have finished
func (s *Service) DontWaitForNotifications() {
	s.notifications = nil
}

// GetNotifications returns the notification channel
func (s *Service) GetNotifications() chan bool {
	return s.notifications
}
