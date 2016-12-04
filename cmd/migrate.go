package cmd

import (
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/util/migrations"
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

	// Run all migrations
	if err := models.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
