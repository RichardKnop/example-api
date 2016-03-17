package accounts

import (
	"errors"
	"fmt"

	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/RichardKnop/recall/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrSuperuserOnlyManually ...
	ErrSuperuserOnlyManually = errors.New("Superusers can only be created manually")
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("User not found")
)

// GetName returns user's full name
func (u *User) GetName() string {
	if u.FirstName.Valid && u.LastName.Valid {
		return fmt.Sprintf("%s %s", u.FirstName.String, u.LastName.String)
	}
	return ""
}

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
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByEmail looks up a user by email and returns it
func (s *Service) FindUserByEmail(email string) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Joins("inner join oauth_users on oauth_users.id = account_users.oauth_user_id").
		Where("oauth_users.username = ?", email).Preload("Account.OauthClient").
		Preload("OauthUser").Preload("Role").First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
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
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByFacebookID looks up a user by a Facebook ID and returns it
func (s *Service) FindUserByFacebookID(facebookID string) (*User, error) {
	// Fetch the user from the database
	user := new(User)
	notFound := s.db.Where("facebook_id = ?", facebookID).
		Preload("Account.OauthClient").Preload("OauthUser").Preload("Role").
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateUser creates a new oauth user and a new account user
func (s *Service) CreateUser(account *Account, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	// Begin a transaction
	tx := s.db.Begin()

	user, err := s.createUserCommon(
		tx,
		account,
		userRequest,
		"",    // facebook ID
		false, // confirmed
	)
	if err != nil {
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

// CreateUserTx creates a new oauth user and a new account user in a transaction
func (s *Service) CreateUserTx(tx *gorm.DB, account *Account, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	return s.createUserCommon(tx, account, userRequest, "", false)
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

// CreateFacebookUser creates a new user with facebook ID
func (s *Service) CreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error) {
	// Superusers can only be created manually
	if userRequest.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	// Begin a transaction
	tx := s.db.Begin()

	user, err := s.createUserCommon(
		tx,
		account,
		userRequest,
		facebookID,
		true, // confirmed
	)
	if err != nil {
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

// CreateSuperuser creates a new superuser account
func (s *Service) CreateSuperuser(account *Account, email, password string) (*User, error) {
	// Begin a transaction
	tx := s.db.Begin()

	userRequest := &UserRequest{
		Email:    email,
		Password: password,
		Role:     roles.Superuser,
	}
	user, err := s.createUserCommon(
		tx,
		account,
		userRequest,
		"",   // facebook ID
		true, // confirmed
	)
	if err != nil {
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

func (s *Service) createUserCommon(db *gorm.DB, account *Account, userRequest *UserRequest, facebookID string, confirmed bool) (*User, error) {
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

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		db,
		userRequest.Email,
		userRequest.Password,
	)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := newUser(
		account,
		oauthUser,
		role,
		facebookID,
		userRequest.FirstName,
		userRequest.LastName,
		confirmed,
	)

	// Save the user to the database
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
