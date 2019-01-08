package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/pidfile"
	"log"
	"os"
	"strconv"
)

type addressVar []string

var (
	serverConfig    = &HTTPConfig{}
	interval        int
	addressesToPing addressVar
	errorRate       float64
	pidFilePath     string
)

// HTTPConfig contains necessary config parameters for a TLS or plain HTTP server
type HTTPConfig struct {
	port       int
	certFile   string
	keyFile    string
	caCertFile string
}

func (c *HTTPConfig) isSecure() bool {
	return c.certFile != "" && c.keyFile != ""
}

func (a *addressVar) String() string {
	return fmt.Sprint(*a)
}

func (a *addressVar) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func init() {
	flag.IntVar(&interval, "interval", 10, "The interval between two pings in seconds")
	flag.Var(&addressesToPing, "address", "The address which should be pinged. Format (http|https)://<IP>:<port>")
	flag.Float64Var(&errorRate, "error-rate", 0.0, "error rate in percent")
	flag.StringVar(&pidFilePath, "pidfile-path", "", "Path to write a PID file to")
	flag.StringVar(&serverConfig.certFile, "cert-path", "", "Path to certificate for https server")
	flag.StringVar(&serverConfig.keyFile, "key-path", "", "Path to key for https server")
	flag.StringVar(&serverConfig.caCertFile, "ca-cert-path", "", "Path to custom ca to trust")

	portFromEnv := os.Getenv("PORT")
	defaultPort := 8080
	if portFromEnv != "" {
		defaultPort, _ = strconv.Atoi(portFromEnv)
	}
	flag.IntVar(&serverConfig.port, "listenPort", defaultPort, "The port to listen on")
}

func main() {
	flag.Parse()

	writePIDFile()

	verifyHTTPSConfig()

	pingServer := NewPingServer(serverConfig, errorRate)
	go pingServer.Start()

	pingClient := NewPingClient(serverConfig, addressesToPing)
	go pingClient.Start()
}

func writePIDFile() {
	if pidFilePath != "" {
		pidfile.SetPidfilePath(pidFilePath)
		err := pidfile.Write()
		if err != nil {
			log.Fatalf("Couldn't write Pidfile at path %s: %s", pidFilePath, err)
		}
	}
}

func verifyHTTPSConfig() {
	if serverConfig.certFile != "" && serverConfig.keyFile == "" {
		log.Fatalf("Specify either certificte and key or none")
	}
	if serverConfig.certFile == "" && serverConfig.keyFile != "" {
		log.Fatalf("Specify either certificte and key or none")
	}
}
