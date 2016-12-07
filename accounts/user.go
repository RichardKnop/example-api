package accounts

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RichardKnop/example-api/logger"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/util"
	"github.com/jinzhu/gorm"
)

var (
	// ErrSuperuserOnlyManually ...
	ErrSuperuserOnlyManually = errors.New("Superusers can only be created manually")
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("User not found")
)

// FindUserByOauthUserID looks up a user by oauth user ID and returns it
func (s *Service) FindUserByOauthUserID(oauthUserID uint) (*models.User, error) {
	// Fetch the user from the database
	user := new(models.User)
	notFound := models.UserPreload(s.db).Where(models.User{
		OauthUserID: util.PositiveIntOrNull(int64(oauthUserID)),
	}).First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByEmail looks up a user by email and returns it
func (s *Service) FindUserByEmail(email string) (*models.User, error) {
	// Fetch the user from the database
	user := new(models.User)
	notFound := models.UserPreload(s.db).
		Joins("inner join oauth_users on oauth_users.id = account_users.oauth_user_id").
		Where("oauth_users.username = LOWER(?)", email).First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByID looks up a user by ID and returns it
func (s *Service) FindUserByID(userID uint) (*models.User, error) {
	// Fetch the user from the database
	user := new(models.User)
	notFound := models.UserPreload(s.db).First(user, userID).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// FindUserByFacebookID looks up a user by a Facebook ID and returns it
func (s *Service) FindUserByFacebookID(facebookID string) (*models.User, error) {
	// Fetch the user from the database
	user := new(models.User)
	notFound := models.UserPreload(s.db).Where("facebook_id = ?", facebookID).
		First(user).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CreateUser creates a new oauth user and a new account user
func (s *Service) CreateUser(account *models.Account, data *UserRequest) (*models.User, error) {
	// Superusers can only be created manually
	if data.Role == roles.Superuser {
		return nil, ErrSuperuserOnlyManually
	}

	// Begin a transaction
	tx := s.db.Begin()

	user, err := s.createUserCommon(
		tx,
		account,
		data,
		"",    // facebook ID
		false, // confirmed
	)
	if err != nil {
		tx.Rollback() // rollback the transaction
		return nil, err
	}

	// Create a new confirmation
	confirmation, err := models.NewConfirmation(user, s.cnf.AppSpecific.ConfirmationLifetime)
	if err != nil {
		return nil, err
	}

	// Save the confirmation to the database
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
			logger.ERROR.Print(err)
		}
	}()

	return user, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(user *models.User, data *UserRequest) error {
	// Is this a request to change user password?
	if data.NewPassword != "" {
		// Verify the user submitted current password
		_, err := s.oauthService.AuthUser(user.OauthUser.Username, data.Password)
		if err != nil {
			return err
		}

		// Set the new password
		return s.oauthService.SetPassword(user.OauthUser, data.NewPassword)
	}

	// Update user metadata
	return s.db.Model(user).UpdateColumns(map[string]interface{}{
		"first_name": util.StringOrNull(data.FirstName),
		"last_name":  util.StringOrNull(data.LastName),
		"updated_at": time.Now(),
	}).Error
}

// GetOrCreateFacebookUser either returns an existing user
// or updates an existing email user with facebook ID or creates a new user
func (s *Service) GetOrCreateFacebookUser(account *models.Account, facebookID string, userRequest *UserRequest) (*models.User, error) {
	var (
		user       *models.User
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
			err = tx.Model(user.OauthUser).UpdateColumns(models.OauthUser{
				Username: userRequest.Email,
				Model:    gorm.Model{UpdatedAt: time.Now()},
			}).Error
			if err != nil {
				tx.Rollback() // rollback the transaction
				return nil, err
			}
		}

		// Set the facebook ID, first name, last name, picture
		err = tx.Model(user).UpdateColumns(models.User{
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
func (s *Service) CreateSuperuser(account *models.Account, email, password string) (*models.User, error) {
	// Begin a transaction
	tx := s.db.Begin()

	data := &UserRequest{
		Email:    email,
		Password: password,
		Role:     roles.Superuser,
	}
	user, err := s.createUserCommon(
		tx,
		account,
		data,
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

// PaginatedUsersCount returns a total count of users
func (s *Service) PaginatedUsersCount() (int, error) {
	var count int
	if err := s.paginatedUsersQuery().Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindPaginatedUsers returns paginated user records
func (s *Service) FindPaginatedUsers(offset, limit int, sorts map[string]string) ([]*models.User, error) {
	var users []*models.User

	// Get the pagination query
	usersQuery := s.paginatedUsersQuery()

	// Sort thre results
	var orderBy []string
	if len(sorts) == 0 {
		orderBy = append(orderBy, "id") // order by ID by default
	} else {
		for field, direction := range sorts {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", field, direction))
		}
	}

	// Retrieve paginated results from the database
	err := models.UserPreload(usersQuery).Offset(offset).Limit(limit).
		Order(strings.Join(orderBy, ",")).Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

// paginatedUsersQuery returns a db query for paginated users
func (s *Service) paginatedUsersQuery() *gorm.DB {
	// Basic query
	usersQuery := s.db.Model(new(models.User))

	return usersQuery
}

func (s *Service) createUserCommon(db *gorm.DB, account *models.Account, data *UserRequest, facebookID string, confirmed bool) (*models.User, error) {
	// Check if email is already taken
	if s.GetOauthService().UserExists(data.Email) {
		return nil, oauth.ErrUsernameTaken
	}

	// If a role is not defined in the user request,
	// fall back to the user role
	if data.Role == "" {
		data.Role = roles.User
	}

	// Create a new oauth user
	oauthUser, err := s.GetOauthService().CreateUserTx(
		db,
		data.Role,
		data.Email,
		data.Password,
	)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user, err := models.NewUser(
		account,
		oauthUser,
		facebookID,
		data.FirstName,
		data.LastName,
		data.Picture,
		confirmed,
	)
	if err != nil {
		return nil, err
	}

	// Save the user to the database
	if err = db.Create(user).Error; err != nil {
		return nil, err
	}

	// Assign related objects
	user.Account = account
	user.OauthUser = oauthUser

	// Update the meta user ID field
	err = db.Model(oauthUser).UpdateColumn(models.OauthUser{MetaUserID: user.ID}).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
