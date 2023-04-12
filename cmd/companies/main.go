package main

import (
	"os"

	"github.com/018bf/companies"
	"github.com/018bf/companies/internal/containers"
	"github.com/urfave/cli/v2"
)

var (
	configPath = ""
)

func main() {
	app := &cli.App{
		Name:    companies.Name,
		Usage:   "service",
		Version: companies.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				EnvVars:     []string{"COMPANIES_CONFIG_PATH"},
				TakesFile:   true,
				Value:       configPath,
				Destination: &configPath,
				HasBeenSet:  false,
			},
		},
		Action: runApp,
		Commands: []*cli.Command{
			{
				Name:      "migrate",
				Usage:     "Run migrations",
				Action:    runMigrations,
				ArgsUsage: "",
			},
			{
				Name:      "grpc",
				Usage:     "Run gRPC server",
				Action:    runGRPC,
				ArgsUsage: "",
			},
			{
				Name:      "rest",
				Usage:     "Run rest server",
				Action:    runREST,
				ArgsUsage: "",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// runApp - run app
func runApp(context *cli.Context) error {
	app := containers.NewRESTContainer(configPath)
	app.Run()
	return nil
}

// runGRPC - run grpc api
func runGRPC(context *cli.Context) error {
	app := containers.NewGRPCContainer(configPath)
	app.Run()
	return nil
}

// runREST - run REST api
func runREST(context *cli.Context) error {
	app := containers.NewRESTContainer(configPath)
	app.Run()
	return nil
}

// runMigrations - migrate database
func runMigrations(context *cli.Context) error {
	app := containers.NewMigrateContainer(configPath)
	app.Run()
	return nil
}
