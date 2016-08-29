package accounts

import (
	"errors"
	"fmt"
	"time"

	"github.com/RichardKnop/example-api/accounts/roles"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
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
	notFound := UserPreload(s.db).Where(User{
		OauthUserID: util.PositiveIntOrNull(int64(oauthUserID)),
	}).First(user).RecordNotFound()

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
	notFound := UserPreload(s.db).
		Joins("inner join oauth_users on oauth_users.id = account_users.oauth_user_id").
		Where("oauth_users.username = LOWER(?)", email).First(user).RecordNotFound()

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
	notFound := UserPreload(s.db).First(user, userID).RecordNotFound()

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
	notFound := UserPreload(s.db).Where("facebook_id = ?", facebookID).
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

	// Create a new confirmation
	confirmation := NewConfirmation(user)
	if err := tx.Create(confirmation).Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Assign related object
	confirmation.User = user

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Send confirmation email
	go func() {
		if err := s.sendConfirmationEmail(confirmation); err != nil {
			logger.Error(err)
		}
	}()

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
	// Is this a request to change user password?
	if userRequest.NewPassword != "" {
		// Verify the user submitted current password
		_, err := s.oauthService.AuthUser(
			user.OauthUser.Username,
			userRequest.Password,
		)
		if err != nil {
			return err
		}

		// Set the new password
		return s.oauthService.SetPassword(
			user.OauthUser,
			userRequest.NewPassword,
		)
	}

	// Update user metadata
	return s.db.Model(user).UpdateColumns(map[string]interface{}{
		"first_name": util.StringOrNull(userRequest.FirstName),
		"last_name":  util.StringOrNull(userRequest.LastName),
		"updated_at": time.Now(),
	}).Error
}

// GetOrCreateFacebookUser either returns an existing user
// or updates an existing email user with facebook ID or creates a new user
func (s *Service) GetOrCreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error) {
	var (
		user       *User
		err        error
		userExists bool
	)

	// Does a user with this facebook ID already exist?
	user, err = s.FindUserByFacebookID(facebookID)
	// User with this facebook ID alraedy exists
	if err == nil {
		userExists = true
	}

	if userExists == false {
		// Does a user with this email already exist?
		user, err = s.FindUserByEmail(userRequest.Email)
		// User with this email already exists
		if err == nil {
			userExists = true
		}
	}

	// Begin a transaction
	tx := s.db.Begin()

	// User already exists, update the record and return
	if userExists {
		if userRequest.Email != user.OauthUser.Username {
			// Update the email if it changed (should not happen)
			err = tx.Model(user.OauthUser).UpdateColumns(oauth.User{
				Username: userRequest.Email,
				Model:    gorm.Model{UpdatedAt: time.Now()},
			}).Error
			if err != nil {
				tx.Rollback() // rollback the transaction
				return nil, err
			}
		}

		// Set the facebook ID, first name, last name, picture
		err = tx.Model(user).UpdateColumns(User{
			FacebookID: util.StringOrNull(facebookID),
			FirstName:  util.StringOrNull(userRequest.FirstName),
			LastName:   util.StringOrNull(userRequest.LastName),
			Picture:    util.StringOrNull(userRequest.Picture),
			Confirmed:  true,
			Model:      gorm.Model{UpdatedAt: time.Now()},
		}).Error
		if err != nil {
			tx.Rollback() // rollback the transaction
			return nil, err
		}

		// Commit the transaction
		if err = tx.Commit().Error; err != nil {
			tx.Rollback() // rollback the transaction
			return nil, err
		}

		return user, nil
	}

	// Facebook registration only creates regular users
	userRequest.Role = roles.User

	user, err = s.createUserCommon(
		tx,
		account,
		userRequest,
		facebookID, // facebook ID
		true,       // confirmed
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
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
	// Check if email is already taken
	if s.GetOauthService().UserExists(userRequest.Email) {
		return nil, oauth.ErrUsernameTaken
	}

	// If a role is not defined in the user request,
	// fall back to the user role
	if userRequest.Role == "" {
		userRequest.Role = roles.User
	}

	// Fetch the role object
	role, err := s.FindRoleByID(userRequest.Role)
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
	user := NewUser(
		account,
		oauthUser,
		role,
		facebookID,
		confirmed,
		&UserRequest{
			FirstName: userRequest.FirstName,
			LastName:  userRequest.LastName,
			Picture:   userRequest.Picture,
		},
	)

	// Save the user to the database
	if err = db.Create(user).Error; err != nil {
		return nil, err
	}

	// Assign related objects
	user.Account = account
	user.OauthUser = oauthUser
	user.Role = role

	// Update the meta user ID field
	err = db.Model(oauthUser).UpdateColumn(oauth.User{MetaUserID: user.ID}).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) sendConfirmationEmail(confirmation *Confirmation) error {
	confirmationEmail, err := s.emailFactory.NewConfirmationEmail(confirmation)
	if err != nil {
		return fmt.Errorf("New confirmation email error: %s", err)
	}

	// Try to send the confirmation email
	if err := s.emailService.Send(confirmationEmail); err != nil {
		return fmt.Errorf("Send email error: %s", err)
	}

	// If the email was sent successfully, update the email_sent flag
	now := time.Now()
	return s.db.Model(confirmation).UpdateColumns(Confirmation{
		EmailSent:   true,
		EmailSentAt: util.TimeOrNull(&now),
		Model:       gorm.Model{UpdatedAt: now},
	}).Error
}
