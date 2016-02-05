package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/oauth"
)

// CreateSuperadmin creates a new superadmin
func CreateSuperadmin() error {
	cnf, db, err := initConfigDB(true, false)
	defer db.Close()
	if err != nil {
		return err
	}

	// Initialise the accounts service
	accountsService := accounts.NewService(cnf, db, oauth.NewService(cnf, db))

	// Fetch the account (assume all superadmins belong to the first account)
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

	// Create a new account
	_, err = accountsService.CreateSuperadmin(
		account,
		strings.TrimRight(email, "\n"),
		strings.TrimRight(password, "\n"),
	)

	return err
}
