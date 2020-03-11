package cmd

import (
	"fmt"
	"log"
	"runtime"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ranskills/c2m/broker"
	"github.com/ranskills/c2m/setting"
)

// CreateListenAction Create listen action
func CreateListenAction(cfg setting.Config) func(args []string, options map[string]string) int {

	return func(args []string, options map[string]string) int {
		topic, _ := options["topic"]

		if topic == "" {
			topic = cfg.Mqtt.Topic
		}
		log.Println("Topic", topic)

		cfg.Mqtt.ClientID += "-sub"
		client := broker.GetClient(cfg)

		go func() {
			messageHandler := func(client mqtt.Client, msg mqtt.Message) {
				go func() {
					fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
				}()
			}
			if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
				log.Println(token.Error())
			}
		}()

		runtime.Goexit()

		return 0
	}
}
