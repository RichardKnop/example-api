package accounts

import (
	"fmt"

	"github.com/RichardKnop/recall/migrations"
	"github.com/jinzhu/gorm"
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	if err := migrate0001(db); err != nil {
		return err
	}

	return nil
}

// Migrate0001 creates accounts schema
func migrate0001(db *gorm.DB) error {
	migrationName := "accounts_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

	var err error

	// Create account_accounts table
	if err := db.CreateTable(new(Account)).Error; err != nil {
		return fmt.Errorf("Error creating account_accounts table: %s", err)
	}

	// Create account_roles table
	if err := db.CreateTable(new(Role)).Error; err != nil {
		return fmt.Errorf("Error creating account_roles table: %s", err)
	}

	// Create account_users table
	if err := db.CreateTable(new(User)).Error; err != nil {
		return fmt.Errorf("Error creating account_users table: %s", err)
	}

	// Add foreign key on account_accounts.oauth_client_id
	err = db.Model(new(Account)).AddForeignKey(
		"oauth_client_id",
		"oauth_clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_accounts.oauth_client_id for oauth_clients(id): %s", err)
	}

	// Add foreign key on account_users.account_id
	err = db.Model(new(User)).AddForeignKey(
		"account_id",
		"account_accounts(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.account_id for account_accounts(id): %s", err)
	}

	// Add foreign key on account_users.oauth_user_id
	err = db.Model(new(User)).AddForeignKey(
		"oauth_user_id",
		"oauth_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.oauth_user_id for oauth_users(id): %s", err)
	}

	// Add foreign key on account_users.role_id
	err = db.Model(new(User)).AddForeignKey(
		"role_id",
		"account_roles(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"account_users.role_id for account_roles(id): %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
