package accounts

import (
	"fmt"

	"github.com/RichardKnop/example-api/migrations"
	"github.com/jinzhu/gorm"
)

var (
	list = []migrations.MigrationStage{
		{"accounts_initial", migrate0001},
	}
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	return migrations.Migrate(db, list)
}

func migrate0001(db *gorm.DB, name string) error {
	var err error

	// Create tables
	if err := db.CreateTable(new(Account)).Error; err != nil {
		return fmt.Errorf("Error creating account_accounts table: %s", err)
	}
	if err := db.CreateTable(new(Role)).Error; err != nil {
		return fmt.Errorf("Error creating account_roles table: %s", err)
	}
	if err := db.CreateTable(new(User)).Error; err != nil {
		return fmt.Errorf("Error creating account_users table: %s", err)
	}
	if err := db.CreateTable(new(Confirmation)).Error; err != nil {
		return fmt.Errorf("Error creating account_confirmations table: %s", err)
	}
	if err := db.CreateTable(new(Invitation)).Error; err != nil {
		return fmt.Errorf("Error creating account_invitations table: %s", err)
	}
	if err := db.CreateTable(new(PasswordReset)).Error; err != nil {
		return fmt.Errorf("Error creating account_password_resets table: %s", err)
	}

	// Add foreign keys
	err = db.Model(new(Account)).AddForeignKey(
		"oauth_client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_accounts.oauth_client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(User)).AddForeignKey(
		"account_id", "account_accounts(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.account_id for account_accounts(id): %s", err)
	}
	err = db.Model(new(User)).AddForeignKey(
		"oauth_user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.oauth_user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(User)).AddForeignKey(
		"role_id", "account_roles(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.role_id for account_roles(id): %s", err)
	}
	err = db.Model(new(Confirmation)).AddForeignKey(
		"user_id", "account_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_confirmations.user_id for account_users(id): %s", err)
	}
	err = db.Model(new(Invitation)).AddForeignKey(
		"invited_user_id", "account_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_invitations.invited_user_id for account_users(id): %s", err)
	}
	err = db.Model(new(Invitation)).AddForeignKey(
		"invited_by_user_id", "account_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_invitations.invited_by_user_id for account_users(id): %s", err)
	}
	err = db.Model(new(PasswordReset)).AddForeignKey(
		"user_id", "account_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_password_resets.user_id for account_users(id): %s", err)
	}

	return nil
}
