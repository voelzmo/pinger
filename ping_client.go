package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// PingClient is a client that can be started
type PingClient interface {
	Start()
}

type pingClient struct {
	addresses []string
	client    *http.Client
}

func (pc *pingClient) Start() {
	c := time.Tick(time.Duration(interval) * time.Second)
	for range c {
		for _, address := range pc.addresses {
			pc.pingAddress(address)
		}
	}
}

func (pc *pingClient) pingAddress(address string) {
	response, err := pc.client.Get(address + "/ping")
	if err != nil {
		log.Printf("Couldn't ping %s: %s", address, err)
	} else {
		defer response.Body.Close()
		var pingAnswer string
		json.NewDecoder(response.Body).Decode(&pingAnswer)
		log.Printf("Pinged %s, got response: %s, \"%s\"", address, response.Status, pingAnswer)
	}
}

// NewPingClient creates a new client for pinging PingServers
func NewPingClient(serverConfig *HTTPConfig, addresses []string) PingClient {
	client := http.DefaultClient
	if serverConfig.caCertFile != "" {
		caCert, err := ioutil.ReadFile(serverConfig.caCertFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatal(err)
		}
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
	}
	return &pingClient{addresses, client}
}
