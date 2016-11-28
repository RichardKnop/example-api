package cmd

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/migrations"
	"github.com/RichardKnop/example-api/oauth"
)

// Migrate runs database migrations
func Migrate() error {
	_, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := oauth.MigrateAll(db); err != nil {
		return err
	}

	// Run migrations for the accounts service
	if err := accounts.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
