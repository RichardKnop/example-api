package accounts

import (
	"database/sql"

	"github.com/RichardKnop/recall/oauth"
	"github.com/RichardKnop/recall/util"
	"github.com/jinzhu/gorm"
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
	AccountID   sql.NullInt64 `sql:"index;not null"`
	OauthUserID sql.NullInt64 `sql:"index;not null"`
	RoleID      sql.NullInt64 `sql:"index;not null"`
	Account     *Account
	OauthUser   *oauth.User
	Role        *Role
	FirstName   sql.NullString `sql:"type:varchar(100)"`
	LastName    sql.NullString `sql:"type:varchar(100)"`
}

// TableName specifies table name
func (u *User) TableName() string {
	return "account_users"
}

// newAccount creates new Account instance
func newAccount(oauthClient *oauth.Client, name, description string) *Account {
	oauthClientID := util.IntOrNull(int64(oauthClient.ID))
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
func newUser(account *Account, oauthUser *oauth.User, role *Role, firstName, lastName string) *User {
	accountID := util.IntOrNull(int64(account.ID))
	oauthUserID := util.IntOrNull(int64(oauthUser.ID))
	roleID := util.IntOrNull(int64(role.ID))
	user := &User{
		AccountID:   accountID,
		OauthUserID: oauthUserID,
		RoleID:      roleID,
		FirstName:   util.StringOrNull(firstName),
		LastName:    util.StringOrNull(lastName),
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
