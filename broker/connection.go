package broker

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ranskills/c2m/setting"
)

// GetClient Connects to the broker returns a client
func GetClient(cfg setting.Config) mqtt.Client {
	uri := getUri(cfg)

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
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetCleanSession(false)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connection Status: ", client.IsConnected())
	})

	//tlsconfig := NewTLSConfig()
	//opts.SetTLSConfig(tlsconfig)

	return opts
}

func listen(cfg setting.Config) {
	client := connect(cfg, getUri(cfg))
	client.Subscribe(cfg.Mqtt.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)
	//
	certpool := x509.NewCertPool()
	//pemCerts, err := ioutil.ReadFile("/Users/ransfordokpoti/dev/bellwethercoffee/mp/certs/ca.crt")
	pemCerts, err := ioutil.ReadFile("./certs/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}
	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("./certs/device001.crt", "./certs/device001.key")
	if err != nil {
		panic(err)
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}
