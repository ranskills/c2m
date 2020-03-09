package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ranskills/c2m/cmd"
	"github.com/ranskills/c2m/setting"
	"github.com/ranskills/c2m/util"
	"github.com/teris-io/cli"
	"log"
	"os"
)

func main() {
	var cfg setting.Config

	configPath := util.GetConfigOptionValue(os.Args)

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	version := cli.NewCommand("version", "print mp version").
		WithAction(cmd.VersionHandler)

	health := cli.NewCommand("health", "reports the status of the message broker").
		WithAction(cmd.CreateHealthAction(cfg))

	dryRun := cli.NewCommand("dry-run", "Processes the files and dumps the output to the console without publishing to the broker").
		WithAction(cmd.CreateDryRunAction(cfg))

	run := cli.NewCommand("run", "Processes the files and publishes to the message broker").
		WithAction(cmd.CreateRunAction(cfg))

	app := cli.New("A customizable tool for publishing the contents of CSV files to a MQTT message broker").
		WithOption(cli.NewOption("config", "full path to a configuration file").WithChar('c').WithType(cli.TypeString)).
		WithOption(cli.NewOption("src", "path to the directory containing files to be published").WithChar('s').WithType(cli.TypeString)).
		WithOption(cli.NewOption("pretty", "to have JSON data properly formatted").WithChar('p').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("watch", "to actively watch for new files in the source directory").WithChar('w').WithType(cli.TypeBool)).
		WithCommand(version).
		WithCommand(health).
		WithCommand(dryRun).
		WithCommand(run)

	os.Exit(app.Run(os.Args, os.Stdout))
}
