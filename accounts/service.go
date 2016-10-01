package accounts

import (
	"github.com/jinzhu/gorm"
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/oauth"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	oauthService oauth.ServiceInterface
	emailService email.ServiceInterface
	emailFactory EmailFactoryInterface
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
