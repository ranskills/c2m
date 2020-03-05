package cmd

import (
	"fmt"
	"strings"

	"github.com/ranskills/mp/broker"
	"github.com/ranskills/mp/setting"
)

// CreateHealthAction Creates the action handler for health
func CreateHealthAction(cfg setting.Config) func(args []string, options map[string]string) int {

	return func(args []string, options map[string]string) int {

		client := broker.GetClient(cfg)

		status := "Ok"
		if !client.IsConnected() {
			status = "Failed"
		}

		fmt.Println("Health Checks")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("\nBroker @%s Status: %s\n\n", cfg.Mqtt.Host, status)

		return 0
	}
}
