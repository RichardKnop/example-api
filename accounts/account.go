package accounts

import (
	"errors"

	"github.com/RichardKnop/recall/util"
)

var (
	errAccountNotFound = errors.New("Account not found")
)

// FindAccountByOauthClientID looks up an account by oauth client ID and returns it
func (s *Service) FindAccountByOauthClientID(oauthClientID uint) (*Account, error) {
	// Fetch the client from the database
	account := new(Account)
	notFound := s.db.Where(Account{
		OauthClientID: util.IntOrNull(int64(oauthClientID)),
	}).Preload("OauthClient").First(account).RecordNotFound()

	// Not found
	if notFound {
		return nil, errAccountNotFound
	}

	return account, nil
}

// FindAccountByID looks up an account by ID and returns it
func (s *Service) FindAccountByID(accountID uint) (*Account, error) {
	// Fetch the client from the database
	account := new(Account)
	notFound := s.db.Preload("OauthClient").
		First(account, accountID).RecordNotFound()

	// Not found
	if notFound {
		return nil, errAccountNotFound
	}

	return account, nil
}

// CreateAccount creates a new account
func (s *Service) CreateAccount(name, description, key, secret, redirectURI string) (*Account, error) {
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
	account := newAccount(oauthClient, name, description)

	// Save the account to the database
	if err := tx.Create(account).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return account, nil
}
