package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/facebookgo/pidfile"
	"github.com/voelzmo/pinger/graphite"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type addressVar []string

var (
	config           *httpConfig = &httpConfig{}
	interval         int
	addressesToPing  addressVar
	listenPort       int
	errorRate        float64
	randomGenerator  *rand.Rand
	pidFilePath      string
	graphiteEndpoint string
	metricsPrefix    string

)

type httpConfig struct {
	port       int
	certFile   string
	keyFile    string
	caCertFile string
}

func (c *httpConfig) isSecure() bool {
	return c.certFile != "" && c.keyFile != ""
}

func (a *addressVar) String() string {
	return fmt.Sprint(*a)
}

func (a *addressVar) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func pingAddress(address string) {
	client, err := createClient()
	if err != nil {
		log.Printf("Couldn't create client: %s", err)
	}
	response, err := client.Get(address + "/ping")

	if err != nil {
		log.Printf("Couldn't ping %s: %s", address, err)
	} else {
		defer response.Body.Close()
		var pingAnswer string
		json.NewDecoder(response.Body).Decode(&pingAnswer)
		log.Printf("Pinged %s, got response: %s, \"%s\"", address, response.Status, pingAnswer)
	}
}

func createClient() (*http.Client, error) {
	client := http.DefaultClient
	if config.caCertFile != "" {
		caCert, err := ioutil.ReadFile(config.caCertFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
	}
	return client, nil
}

func startHTTPServer() {
	var pingMetric *graphite.Metric
	if graphiteEndpoint != "" {
		sender, err := graphite.NewGraphiteSender(graphiteEndpoint)
		if err == nil {
			pingMetric = graphite.NewMetric(metricsPrefix, 10.0, sender)
		}
	}
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if pingMetric != nil {
			pingMetric.Increment()
		}
		f := randomGenerator.Float64()
		if f < errorRate {
			http.Error(w, "no wai!", http.StatusTeapot)
			return
		}
		json.NewEncoder(w).Encode("pong")
		log.Printf("Got pinged by %s, sent a 'pong'!", r.RemoteAddr)
	})
	log.Fatal(listenAndServe(config))
}

func listenAndServe(config *httpConfig) error {
	address := fmt.Sprintf(":%v", config.port)
	if !config.isSecure() {
		return http.ListenAndServe(address, nil)
	} else {
		return http.ListenAndServeTLS(address, config.certFile, config.keyFile, nil)
	}
}

func init() {
	flag.IntVar(&interval, "interval", 10, "The interval between two pings in seconds")
	flag.Var(&addressesToPing, "address", "The address which should be pinged. Format (http|https)://<IP>:<port>")
	flag.Float64Var(&errorRate, "error-rate", 0.0, "error rate in percent")
	flag.StringVar(&pidFilePath, "pidfile-path", "", "Path to write a PID file to")
	flag.StringVar(&config.certFile, "cert-path", "", "Path to certificate for https server")
	flag.StringVar(&config.keyFile, "key-path", "", "Path to key for https server")
	flag.StringVar(&config.caCertFile, "ca-cert-path", "", "Path to custom ca to trust")
	flag.StringVar(&graphiteEndpoint, "graphite-endpoint", "", "Where to write metrics to, format is <host>:<port>")
	flag.StringVar(&metricsPrefix, "metrics-prefix", "pinger", "Prefix for reporting metrics")

	portFromEnv := os.Getenv("PORT")
	defaultPort := 8080
	if portFromEnv != "" {
		defaultPort, _ = strconv.Atoi(portFromEnv)
	}
	flag.IntVar(&config.port, "listenPort", defaultPort, "The port to listen on")
}

func main() {
	flag.Parse()
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

	if pidFilePath != "" {
		pidfile.SetPidfilePath(pidFilePath)
		err := pidfile.Write()
		if err != nil {
			log.Fatalf("Couldn't write Pidfile at path %s: %s", pidFilePath, err)
		}
	}

	verifyHttpsConfig()

	go startHTTPServer()
	c := time.Tick(time.Duration(interval) * time.Second)
	for range c {
		for _, address := range addressesToPing {
			pingAddress(address)
		}
	}
}

func verifyHttpsConfig() {
	if config.certFile != "" && config.keyFile == "" {
		log.Fatalf("Specify either certificte and key or none")
	}
	if config.certFile == "" && config.keyFile != "" {
		log.Fatalf("Specify either certificte and key or none")
	}
}
