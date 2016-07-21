package main

import (
	"log"
	"os"

	"github.com/RichardKnop/example-api/commands"
	"github.com/urfave/cli"
)

var (
	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "example-api"
	cliApp.Usage = "example-api"
	cliApp.Author = "Richard Knop"
	cliApp.Email = "risoknop@gmail.com"
	cliApp.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) error {
				return commands.Migrate()
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) error {
				return commands.LoadData(c.Args())
			},
		},
		{
			Name:  "createaccount",
			Usage: "create new account",
			Action: func(c *cli.Context) error {
				return commands.CreateAccount()
			},
		},
		{
			Name:  "createsuperuser",
			Usage: "create new superuser",
			Action: func(c *cli.Context) error {
				return commands.CreateSuperuser()
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				return commands.RunServer()
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
