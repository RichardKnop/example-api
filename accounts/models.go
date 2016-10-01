package accounts

import (
	"database/sql"
	"time"

	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/uuid"
	"github.com/jinzhu/gorm"
)

// Account represents an extension of Oauth 2.0 client,
// can be used to group users together
type Account struct {
	gorm.Model
	OauthClientID sql.NullInt64 `sql:"index;not null"`
	OauthClient   *oauth.Client
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
	OauthUser   *oauth.User
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

// EmailTokenModel is an abstract model which can be used for objects from which
// we derive redirect emails (email confirmation, password reset and such)
type EmailTokenModel struct {
	gorm.Model
	Reference   string `sql:"type:varchar(40);unique;not null"`
	EmailSent   bool   `sql:"index;not null"`
	EmailSentAt *time.Time
	ExpiresAt   time.Time `sql:"index;not null"`
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
func NewAccount(oauthClient *oauth.Client, name, description string) (*Account, error) {
	oauthClientID := util.PositiveIntOrNull(int64(oauthClient.ID))
	account := &Account{
		OauthClientID: oauthClientID,
		Name:          name,
		Description:   util.StringOrNull(description),
	}
	return account, nil
}

// NewUser creates new User instance
func NewUser(account *Account, oauthUser *oauth.User, facebookID string, confirmed bool, data *UserRequest) (*User, error) {
	accountID := util.PositiveIntOrNull(int64(account.ID))
	oauthUserID := util.PositiveIntOrNull(int64(oauthUser.ID))
	user := &User{
		AccountID:   accountID,
		OauthUserID: oauthUserID,
		FacebookID:  util.StringOrNull(facebookID),
		FirstName:   util.StringOrNull(data.FirstName),
		LastName:    util.StringOrNull(data.LastName),
		Picture:     util.StringOrNull(data.Picture),
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
