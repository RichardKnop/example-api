package facebook

import (
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/jinzhu/gorm"
)

// Service struct keeps objects to avoid passing them around
type Service struct {
	cnf             *config.Config
	db              *gorm.DB
	accountsService accounts.ServiceInterface
	adapter         AdapterInterface
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, accountsService accounts.ServiceInterface, adapter AdapterInterface) *Service {
	if adapter == nil {
		adapter = NewAdapter(cnf)
	}
	return &Service{
		cnf:             cnf,
		db:              db,
		accountsService: accountsService,
		adapter:         adapter,
	}
}

// GetAccountsService returns accounts.Service instance
func (s *Service) GetAccountsService() accounts.ServiceInterface {
	return s.accountsService
}

// GetAdapter returns Adapter instance
func (s *Service) GetAdapter() AdapterInterface {
	return s.adapter
}
