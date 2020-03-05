package broker

import (
	"fmt"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ranskills/mp/setting"
)

// GetClient Connects to the broker returns a client
func GetClient(cfg setting.Config) mqtt.Client {
	// clientID string, uri *url.URL
	// "mqtt://127.0.0.1:1883"

	uri := getUri(cfg)
	//listen(cfg)
	return connect(cfg, uri)
}

func getUri(cfg setting.Config) *url.URL {
	uri, err := url.Parse(fmt.Sprintf("%s://%s:%s", cfg.Mqtt.Scheme, cfg.Mqtt.Host, cfg.Mqtt.Port))
	if err != nil {
		log.Fatal(err)
	}

	return uri
}

func connect(cfg setting.Config, uri *url.URL) mqtt.Client {
	opts := createClientOptions(cfg, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return client
}

func createClientOptions(cfg setting.Config, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	serverConfig := cfg.Mqtt

	scheme := serverConfig.Scheme
	opts.AddBroker(fmt.Sprintf("%s://%s", scheme, uri.Host))

	if serverConfig.Username != "" {
		opts.SetUsername(serverConfig.Username)
		opts.SetPassword(serverConfig.Password)
	}
	opts.SetClientID(serverConfig.ClientID)
	opts.SetAutoReconnect(false)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetCleanSession(true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connection Status: ", client.IsConnected())
	})

	return opts
}

func listen(cfg setting.Config) {
	client := connect(cfg, getUri(cfg))
	client.Subscribe(cfg.Mqtt.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}
