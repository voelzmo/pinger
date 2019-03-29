package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/facebookgo/pidfile"
	"github.com/spf13/viper"
)

var (
	configPath string
)

type config struct {
	interval        time.Duration
	addressesToPing []string
	errorRate       float64
	pidFilePath     string
}

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

func init() {
	flag.StringVar(&configPath, "config-path", "", "Absolute path to read a config file from")
}

func main() {
	readConfiguration()

	writePIDFile(viper.GetString("pidfile-path"))

	serverConfig := &HTTPConfig{
		viper.GetInt("listenPort"),
		viper.GetString("cert-path"),
		viper.GetString("key-path"),
		viper.GetString("ca-cert-path"),
	}
	verifyHTTPSConfig(serverConfig)

	pingServer := NewPingServer(serverConfig, viper.GetFloat64("error-rate"))
	go pingServer.Start()

	pingClient := NewPingClient(serverConfig, viper.GetStringSlice("address"), viper.GetDuration("interval"))
	go pingClient.Start()

	select {} // keep main running, as we only have gofuncs above
}

func readConfiguration() {
	flag.Parse()
	if configPath != "" {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file at '%s': %s", configPath, err))
		}
	}
}

func writePIDFile(path string) {
	if path != "" {
		pidfile.SetPidfilePath(path)
		err := pidfile.Write()
		if err != nil {
			log.Fatalf("Couldn't write Pidfile at path '%s': %s", path, err)
		}
	}
}

func verifyHTTPSConfig(sc *HTTPConfig) {
	if sc.certFile != "" && sc.keyFile == "" {
		log.Fatal("Specify either certificate and key or none")
	}
	if sc.certFile == "" && sc.keyFile != "" {
		log.Fatal("Specify either certificate and key or none")
	}
}
