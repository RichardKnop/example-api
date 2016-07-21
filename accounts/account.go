package accounts

import (
	"errors"

	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrAccountNotFound ...
	ErrAccountNotFound = errors.New("Account not found")
	// ErrAccountNameTaken ...
	ErrAccountNameTaken = errors.New("Account name taken")
)

// FindAccountByOauthClientID looks up an account by oauth client ID and returns it
func (s *Service) FindAccountByOauthClientID(oauthClientID uint) (*Account, error) {
	// Fetch the client from the database
	account := new(Account)
	notFound := AccountPreload(s.db).Where(Account{
		OauthClientID: util.PositiveIntOrNull(int64(oauthClientID)),
	}).First(account).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

// FindAccountByID looks up an account by ID and returns it
func (s *Service) FindAccountByID(accountID uint) (*Account, error) {
	// Fetch the client from the database
	account := new(Account)
	notFound := AccountPreload(s.db).First(account, accountID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

// FindAccountByName looks up an account by name and returns it
func (s *Service) FindAccountByName(name string) (*Account, error) {
	// Fetch the client from the database
	account := new(Account)
	notFound := AccountPreload(s.db).Where("name = ?", name).
		First(account).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

// CreateAccount creates a new account
func (s *Service) CreateAccount(name, description, key, secret, redirectURI string) (*Account, error) {
	// Check uniqueness of the name
	account, err := s.FindAccountByName(name)
	if account != nil && err == nil {
		return nil, ErrAccountNameTaken
	}

	// Check uniqueness of the key (client ID)
	if s.GetOauthService().ClientExists(key) {
		return nil, oauth.ErrClientIDTaken
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new oauth client
	oauthClient, err := s.GetOauthService().CreateClientTx(
		tx,
		key,
		secret,
		redirectURI,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new account
	account = NewAccount(oauthClient, name, description)

	// Save the account to the database
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign related object
	account.OauthClient = oauthClient

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return account, nil
}

// AccountPreload sets up Gorm preloads for an account object
func AccountPreload(db *gorm.DB) *gorm.DB {
	return AccountPreloadWithPrefix(db, "")
}

// AccountPreloadWithPrefix sets up Gorm preloads for an account object, and prefixes with prefix for nested objects
func AccountPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "OauthClient")
}
