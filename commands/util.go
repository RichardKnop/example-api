package commands

import (
	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/database"
	"github.com/jinzhu/gorm"
)

// initConfigDB loads the configuration and connects to the database
func initConfigDB(mustLoadOnce, keepReloading bool) (*config.Config, *gorm.DB, error) {
	// Config
	cnf := config.NewConfig(mustLoadOnce, keepReloading)

	// Database
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return nil, nil, err
	}

	// Database logging
	db.LogMode(cnf.IsDevelopment)

	return cnf, db, nil
}
