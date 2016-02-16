package accounts

import (
	"errors"

	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/RichardKnop/recall/util"
)

var (
	errSuperuserOnlyManually = errors.New("Superusers can only be created manually")
	errUserNotFound          = errors.New("User not found")
)

// FindUserByOauthUserID looks up a user by oauth user ID and returns it
func (s *Service) FindUserByOauthUserID(oauthUserID uint) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Where(User{
		OauthUserID: util.PositiveIntOrNull(int64(oauthUserID)),
	}).Preload("Account.OauthClient").Preload("OauthUser").Preload("Role").
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, errUserNotFound
	}

	return user, nil
}

// FindUserByID looks up a user by ID and returns it
func (s *Service) FindUserByID(userID uint) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Preload("Account.OauthClient").Preload("OauthUser").
		Preload("Role").First(user, userID).RecordNotFound()

	// Not found
	if notFound {
		return nil, errUserNotFound
	}

	return user, nil
}

// CreateUser creates a new oauth user and a new account user
func (s *Service) CreateUser(account *Account, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, errSuperuserOnlyManually
	}

	// If a role is not defined in the user request,
	// fall back to the user role
	if userRequest.Role == "" {
		userRequest.Role = roles.User
	}

	// Fetch the role object
	role, err := s.findRoleByName(userRequest.Role)
	if err != nil {
		return nil, err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		tx,
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new user
	user := newUser(
		account,
		oauthUser,
		role,
		userRequest.FirstName,
		userRequest.LastName,
	)

	// Save the user to the database
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(user *User, userRequest *UserRequest) error {
	// Update basic metadata
	if err := s.db.Model(user).UpdateColumns(User{
		FirstName: util.StringOrNull(userRequest.FirstName),
		LastName:  util.StringOrNull(userRequest.LastName),
	}).Error; err != nil {
		return err
	}

	return nil
}

// CreateSuperuser creates a new superuser account
func (s *Service) CreateSuperuser(account *Account, email, password string) (*User, error) {
	// Fetch the role object
	role, err := s.findRoleByName(roles.Superuser)
	if err != nil {
		return nil, err
	}

	// Begin a transaction
	tx := s.db.Begin()

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		tx,
		email,
		password,
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new user
	user := newUser(
		account,
		oauthUser,
		role,
		"",
		"",
	)

	// Save the user to the database
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	return user, nil
}
