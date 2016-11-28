package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/oauth"
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

	// Fetch the account (assume all superusers belong to the first account)
	account, err := accountsService.FindAccountByID(uint(1))
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

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
		account,
		strings.TrimRight(email, "\n"),
		strings.TrimRight(password, "\n"),
	)

	return err
}
