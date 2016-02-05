package web

import (
	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/config"
)

// Service struct keeps variables for reuse
type Service struct {
	cnf             *config.Config
	accountsService accounts.ServiceInterface
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, accountsService accounts.ServiceInterface) *Service {
	return &Service{
		cnf:             cnf,
		accountsService: accountsService,
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}
