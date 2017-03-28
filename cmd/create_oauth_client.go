package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/example-api/services/oauth"
)

// CreateOauthClient creates a new OAuth client
func CreateOauthClient() error {
	cnf, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	reader := bufio.NewReader(os.Stdin)

	// OAuth client ID from a user input
	fmt.Print("Client ID: ")
	oauthClientID, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// OAuth client secret from a user input
	fmt.Print("Secret: ")
	oauthClientSecret, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// OAuth client redirect URI from a user input
	fmt.Print("Redirect URI: ")
	oauthClientRedirectURI, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Create a new account
	_, err = oauthService.CreateClient(
		strings.TrimRight(oauthClientID, "\n"),
		strings.TrimRight(oauthClientSecret, "\n"),
		strings.TrimRight(oauthClientRedirectURI, "\n"),
	)

	return err
}
