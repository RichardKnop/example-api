package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/services/email"
	"github.com/RichardKnop/example-api/services/oauth"
)

// CreateSuperuser creates a new superuser
func CreateSuperuser() error {
	cnf, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

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

	reader := bufio.NewReader(os.Stdin)

	// OAuth client ID from a user input
	fmt.Print("Client ID: ")
	oauthClientID, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Fetch the oauth client
	oauthClient, err := oauthService.FindClientByClientID(oauthClientID)
	if err != nil {
		return err
	}

	// Email from a user input
	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Password from a user input
	fmt.Print("Password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Create a new user
	_, err = accountsService.CreateSuperuser(
		oauthClient,
		strings.TrimRight(email, "\n"),
		strings.TrimRight(password, "\n"),
	)

	return err
}
