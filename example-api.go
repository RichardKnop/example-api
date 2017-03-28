package main

import (
	"log"
	"os"

	"github.com/RichardKnop/example-api/cmd"
	"github.com/urfave/cli"
)

var (
	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "Example API CLI"
	cliApp.Usage = "example-api"
	cliApp.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "Run database migrations",
			Action: func(c *cli.Context) error {
				return cmd.Migrate()
			},
		},
		{
			Name:  "loaddata",
			Usage: "Load data from fixture into the database",
			Action: func(c *cli.Context) error {
				return cmd.LoadData(c.Args())
			},
		},
		{
			Name:  "createoauthclient",
			Usage: "Create a new OAuth client",
			Action: func(c *cli.Context) error {
				return cmd.CreateOauthClient()
			},
		},
		{
			Name:  "createsuperuser",
			Usage: "Create a new superuser",
			Action: func(c *cli.Context) error {
				return cmd.CreateSuperuser()
			},
		},
		{
			Name:  "runserver",
			Usage: "Run the web server on port 8080",
			Action: func(c *cli.Context) error {
				return cmd.RunServer()
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
