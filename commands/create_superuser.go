package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/oauth"
)

// CreateSuperuser creates a new superuser
func CreateSuperuser() error {
	cnf, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialise the accounts service
	accountsService := accounts.NewService(cnf, db, oauth.NewService(cnf, db))

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
