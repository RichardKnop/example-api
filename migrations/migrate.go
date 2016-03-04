package migrations

import (
	"github.com/jinzhu/gorm"
)

// MigrateAll runs bootstrap, then all migration functions listed against
// the specified database and logs any errors
func MigrateAll(db *gorm.DB, migrationFunctions []func(*gorm.DB) error) error {

	// Begin a transaction
	tx := db.Begin()

	if err := Bootstrap(tx); err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	for _, m := range migrationFunctions {
		if err := m(tx); err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}
