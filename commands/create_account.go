package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/oauth"
)

// CreateAccount creates a new account
func CreateAccount() error {
	cnf, db, err := initConfigDB(true, false)
	defer db.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	// Account name from a user input
	fmt.Print("Name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Account description from a user input
	fmt.Print("Description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// API key (oauth client ID) from a user input
	fmt.Print("API Key (oauth client ID): ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// OAuth 2.0 client secret from a user input
	fmt.Print("API Secret (oauth client secret): ")
	apiSecret, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// OAuth 2.0 client redirect URI from a user input
	fmt.Print("Redirect URI (oauth): ")
	redirectURI, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Initialise the accounts service
	accountsService := accounts.NewService(cnf, db, oauth.NewService(cnf, db))

	// Create a new account
	_, err = accountsService.CreateAccount(
		strings.TrimRight(name, "\n"),
		strings.TrimRight(description, "\n"),
		strings.TrimRight(apiKey, "\n"),
		strings.TrimRight(apiSecret, "\n"),
		strings.TrimRight(redirectURI, "\n"),
	)

	return err
}
