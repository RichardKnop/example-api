package accounts

import (
	"database/sql"

	"github.com/RichardKnop/recall/oauth"
	"github.com/RichardKnop/recall/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pborman/uuid"
)

// Account ...
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

// Role ...
type Role struct {
	gorm.Model
	Name string `sql:"type:varchar(20);unique;not null"`
}

// TableName specifies table name
func (r *Role) TableName() string {
	return "account_roles"
}

// User ...
type User struct {
	gorm.Model
	AccountID   sql.NullInt64  `sql:"index;not null"`
	OauthUserID sql.NullInt64  `sql:"index;not null"`
	RoleID      sql.NullInt64  `sql:"index;not null"`
	FacebookID  sql.NullString `sql:"type:varchar(60);unique"`
	Account     *Account
	OauthUser   *oauth.User
	Role        *Role
	FirstName   sql.NullString `sql:"type:varchar(100)"`
	LastName    sql.NullString `sql:"type:varchar(100)"`
	Confirmed   bool           `sql:"index;not null"`
}

// TableName specifies table name
func (u *User) TableName() string {
	return "account_users"
}

// Confirmation ...
type Confirmation struct {
	gorm.Model
	UserID      sql.NullInt64 `sql:"index;not null"`
	User        *User
	Reference   string `sql:"type:varchar(40);unique;not null"`
	EmailSent   bool   `sql:"index;not null"`
	EmailSentAt pq.NullTime
}

// TableName specifies table name
func (c *Confirmation) TableName() string {
	return "account_confirmations"
}

// newAccount creates new Account instance
func newAccount(oauthClient *oauth.Client, name, description string) *Account {
	oauthClientID := util.PositiveIntOrNull(int64(oauthClient.ID))
	account := &Account{
		OauthClientID: oauthClientID,
		Name:          name,
		Description:   util.StringOrNull(description),
	}
	if oauthClientID.Valid {
		account.OauthClient = oauthClient
	}
	return account
}

// newUser creates new User instance
func newUser(account *Account, oauthUser *oauth.User, role *Role, facebookID, firstName, lastName string, confirmed bool) *User {
	accountID := util.PositiveIntOrNull(int64(account.ID))
	oauthUserID := util.PositiveIntOrNull(int64(oauthUser.ID))
	roleID := util.PositiveIntOrNull(int64(role.ID))
	user := &User{
		AccountID:   accountID,
		OauthUserID: oauthUserID,
		RoleID:      roleID,
		FacebookID:  util.StringOrNull(facebookID),
		FirstName:   util.StringOrNull(firstName),
		LastName:    util.StringOrNull(lastName),
		Confirmed:   confirmed,
	}
	if accountID.Valid {
		user.Account = account
	}
	if oauthUserID.Valid {
		user.OauthUser = oauthUser
	}
	if roleID.Valid {
		user.Role = role
	}
	return user
}

// newConfirmation creates new Confirmation instance
func newConfirmation(user *User) *Confirmation {
	userID := util.PositiveIntOrNull(int64(user.ID))
	confirmation := &Confirmation{
		UserID:      userID,
		Reference:   uuid.New(),
		EmailSentAt: util.TimeOrNull(nil),
	}
	if userID.Valid {
		confirmation.User = user
	}
	return confirmation
}
