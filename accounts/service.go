package accounts

import (
	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/email"
	"github.com/RichardKnop/recall/oauth"
	"github.com/jinzhu/gorm"
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
