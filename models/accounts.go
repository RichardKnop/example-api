package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/uuid"
	"github.com/jinzhu/gorm"
)

// Account represents an extension of Oauth 2.0 client,
// can be used to group users together
type Account struct {
	gorm.Model
	OauthClientID sql.NullInt64 `sql:"index;not null"`
	OauthClient   *OauthClient
	Name          string         `sql:"type:varchar(100);unique;not null"`
	Description   sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (p *Account) TableName() string {
	return "account_accounts"
}

// User represents a platform user
type User struct {
	gorm.Model
	AccountID   sql.NullInt64 `sql:"index;not null"`
	OauthUserID sql.NullInt64 `sql:"index;not null"`
	Account     *Account
	OauthUser   *OauthUser
	FacebookID  sql.NullString `sql:"type:varchar(60);unique"`
	FirstName   sql.NullString `sql:"type:varchar(100)"`
	LastName    sql.NullString `sql:"type:varchar(100)"`
	Picture     sql.NullString `sql:"type:varchar(255)"`
	Confirmed   bool           `sql:"index;not null"`
}

// TableName specifies table name
func (u *User) TableName() string {
	return "account_users"
}

// GetName returns user's full name
func (u *User) GetName() string {
	if u.FirstName.Valid && u.LastName.Valid {
		return fmt.Sprintf("%s %s", u.FirstName.String, u.LastName.String)
	}
	return ""
}

// Confirmation objects is created when we send user a confirmation email
// It is then fetched when user clicks on the verification link in the email
// so we can verify his/her email
type Confirmation struct {
	EmailTokenModel
	UserID sql.NullInt64 `sql:"index;not null"`
	User   *User
}

// TableName specifies table name
func (c *Confirmation) TableName() string {
	return "account_confirmations"
}

// Invitation is created when user invites another user to the platform.
// We send out an invite email and the invited user can follow the link to
// set a password and finish the sign up process
type Invitation struct {
	EmailTokenModel
	InvitedUserID   sql.NullInt64 `sql:"index;not null"`
	InvitedByUserID sql.NullInt64 `sql:"index;not null"`
	InvitedByUser   *User
	InvitedUser     *User
}

// TableName specifies table name
func (i *Invitation) TableName() string {
	return "account_invitations"
}

// PasswordReset is created when user forgets his/her password and requests
// a new one. We send out an email with a link where user can set a new password.
type PasswordReset struct {
	EmailTokenModel
	UserID sql.NullInt64 `sql:"index;not null"`
	User   *User
}

// TableName specifies table name
func (p *PasswordReset) TableName() string {
	return "account_password_resets"
}

// NewAccount creates new Account instance
func NewAccount(oauthClient *OauthClient, name, description string) (*Account, error) {
	oauthClientID := util.PositiveIntOrNull(int64(oauthClient.ID))
	account := &Account{
		OauthClientID: oauthClientID,
		Name:          name,
		Description:   util.StringOrNull(description),
	}
	return account, nil
}

// NewUser creates new User instance
func NewUser(account *Account, oauthUser *OauthUser, facebookID, firstName, lastName, picture string, confirmed bool) (*User, error) {
	accountID := util.PositiveIntOrNull(int64(account.ID))
	oauthUserID := util.PositiveIntOrNull(int64(oauthUser.ID))
	user := &User{
		AccountID:   accountID,
		OauthUserID: oauthUserID,
		FacebookID:  util.StringOrNull(facebookID),
		FirstName:   util.StringOrNull(firstName),
		LastName:    util.StringOrNull(lastName),
		Picture:     util.StringOrNull(picture),
		Confirmed:   confirmed,
	}
	return user, nil
}

// NewConfirmation creates new Confirmation instance
func NewConfirmation(user *User, expiresIn int) (*Confirmation, error) {
	userID := util.PositiveIntOrNull(int64(user.ID))
	confirmation := &Confirmation{
		EmailTokenModel: EmailTokenModel{
			Reference:   uuid.New(),
			EmailSentAt: nil,
			ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		},
		UserID: userID,
	}
	return confirmation, nil
}

// NewInvitation creates new Invitation instance
func NewInvitation(invitedUser, invitedByUser *User, expiresIn int) (*Invitation, error) {
	invitedUserID := util.PositiveIntOrNull(int64(invitedUser.ID))
	invitedByUserID := util.PositiveIntOrNull(int64(invitedByUser.ID))
	invitation := &Invitation{
		EmailTokenModel: EmailTokenModel{
			Reference:   uuid.New(),
			EmailSentAt: nil,
			ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		},
		InvitedUserID:   invitedUserID,
		InvitedByUserID: invitedByUserID,
	}
	return invitation, nil
}

// NewPasswordReset creates new PasswordReset instance
func NewPasswordReset(user *User, expiresIn int) (*PasswordReset, error) {
	userID := util.PositiveIntOrNull(int64(user.ID))
	passwordReset := &PasswordReset{
		EmailTokenModel: EmailTokenModel{
			Reference:   uuid.New(),
			EmailSentAt: nil,
			ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		},
		UserID: userID,
	}
	return passwordReset, nil
}

// AccountPreload sets up Gorm preloads for an account object
func AccountPreload(db *gorm.DB) *gorm.DB {
	return AccountPreloadWithPrefix(db, "")
}

// AccountPreloadWithPrefix sets up Gorm preloads for an account object, and prefixes with prefix for nested objects
func AccountPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "OauthClient")
}

// UserPreload sets up Gorm preloads for a user object
func UserPreload(db *gorm.DB) *gorm.DB {
	return UserPreloadWithPrefix(db, "")
}

// UserPreloadWithPrefix sets up Gorm preloads for a user object,
// and prefixes with prefix for nested objects
func UserPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "Account.OauthClient").
		Preload(prefix + "OauthUser")
}

// ConfirmationPreload sets up Gorm preloads for a confirmation object
func ConfirmationPreload(db *gorm.DB) *gorm.DB {
	return ConfirmationPreloadWithPrefix(db, "")
}

// ConfirmationPreloadWithPrefix sets up Gorm preloads for a confirmation object,
// and prefixes with prefix for nested objects
func ConfirmationPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "User.OauthUser")
}

// InvitationPreload sets up Gorm preloads for an invitation object
func InvitationPreload(db *gorm.DB) *gorm.DB {
	return InvitationPreloadWithPrefix(db, "")
}

// InvitationPreloadWithPrefix sets up Gorm preloads for an invitation object,
// and prefixes with prefix for nested objects
func InvitationPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "InvitedUser.OauthUser").
		Preload(prefix + "InvitedByUser.OauthUser")
}

// PasswordResetPreload sets up Gorm preloads for a password reset object
func PasswordResetPreload(db *gorm.DB) *gorm.DB {
	return PasswordResetPreloadWithPrefix(db, "")
}

// PasswordResetPreloadWithPrefix sets up Gorm preloads for a password reset object,
// and prefixes with prefix for nested objects
func PasswordResetPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.Preload(prefix + "User.OauthUser")
}
