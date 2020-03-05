package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ranskills/mp/cmd"
	"github.com/ranskills/mp/setting"
	"github.com/teris-io/cli"
	"log"
	"os"
	"strings"
)

var cfg setting.Config

func watchDirectory() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					// processFile(event.Name)
				}
				fmt.Printf("EVENT! %#v\n", event)

			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	dir := "./files"
	if err := watcher.Add(dir); err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Printf("Watching the directory %s for new files\n", dir)
		fmt.Println(strings.Repeat("-", 80))
	}

	<-done
}

func main() {
	//configPath := flag.String("x", "./config.yml", "XXXX")
	//log.Println(*configPath)
	err := cleanenv.ReadConfig("./config.yml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	//client := broker.GetClient(cfg)
	// fmt.Println("Connection Status: ", client.IsConnected())
	//client.Subscribe()
	//client.Publish(cfg.Mqtt.Topic, 0, false, time.Now().String())
	//if token := client.Subscribe(cfg.Mqtt.Topic, 0, nil); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
	//time.Sleep(2 * time.Second)

	version := cli.NewCommand("version", "print mp version").
		WithAction(cmd.VersionHandler)

	health := cli.NewCommand("health", "reports the status of the message broker").
		WithAction(cmd.CreateHealthAction(cfg))

	dryRun := cli.NewCommand("dry-run", "Processes the files and dumps the output to the console without publishing to the broker").
		WithAction(cmd.CreateDryRunAction(cfg))

	run := cli.NewCommand("run", "Processes the files and publishes to the message broker").
		WithAction(func(args []string, options map[string]string) int {
			fmt.Printf("%s\n", "OK")
			return 0
		})

	app := cli.New("A customizable tool for publishing the contents of CSV files to a MQTT message broker").
		WithOption(cli.NewOption("config", "full path to a configuration file").WithChar('c').WithType(cli.TypeString)).
		WithOption(cli.NewOption("src", "path to the directory containing files to be published").WithChar('s').WithType(cli.TypeString)).
		WithOption(cli.NewOption("pretty", "to have JSON data properly formatted").WithChar('p').WithType(cli.TypeBool)).
		WithOption(cli.NewOption("watch", "to actively watch for new files in the source directory").WithChar('p').WithType(cli.TypeBool)).
		WithCommand(version).
		WithCommand(health).
		WithCommand(dryRun).
		WithCommand(run)

	os.Exit(app.Run(os.Args, os.Stdout))
}
