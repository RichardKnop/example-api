package commands

import (
	"time"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/facebook"
	"github.com/RichardKnop/example-api/health"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
)

// RunServer runs the app
func RunServer() error {
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialise the health service
	healthService := health.NewService(db)

	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	// Initialise the email service
	emailService := email.NewService(cnf)

	// Initialise the accounts service
	accountsService := accounts.NewService(
		cnf,
		db,
		oauthService,
		emailService,
		nil, // accounts.EmailFactory
	)

	// Initialise the facebook service
	facebookService := facebook.NewService(
		cnf,
		db,
		accountsService,
		nil, // facebook.Adapter
	)

	// Start a negroni app
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))

	// Create a router instance
	router := mux.NewRouter()

	// Add routes for the health service (healthcheck endpoint)
	health.RegisterRoutes(router, healthService)

	// Add routes for the oauth service (tokens endpoint)
	oauth.RegisterRoutes(router, oauthService)

	// Register routes for the accounts service
	accounts.RegisterRoutes(router, accountsService)

	// Register routes for the facebook service
	facebook.RegisterRoutes(router, facebookService)

	// Set the router
	app.UseHandler(router)

	// Run the server on port 8080, gracefully stop on SIGTERM signal
	graceful.Run(":8080", 10*time.Second, app)

	return nil
}
