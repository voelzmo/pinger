package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/pidfile"
	"github.com/spf13/viper"
	"log"
	"time"
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
	//flag.IntVar(&interval, "interval", 10, "The interval between two pings in seconds")
	//flag.StringArrayVar(&addressesToPing, "address", []string{}, "The address which should be pinged. Format (http|https)://<IP>:<port>")
	//flag.Float64Var(&errorRate, "error-rate", 0.0, "error rate in percent")
	//flag.StringVar(&pidFilePath, "pidfile-path", "", "Path to write a PID file to")
	flag.StringVar(&configPath, "config-path", "", "Absolute path to read a config file from")
	//flag.StringVar(&serverConfig.certFile, "cert-path", "", "Path to certificate for https server")
	//flag.StringVar(&serverConfig.keyFile, "key-path", "", "Path to key for https server")
	//flag.StringVar(&serverConfig.caCertFile, "ca-cert-path", "", "Path to custom ca to trust")

	//portFromEnv := os.Getenv("PORT")
	//defaultPort := 8080
	//if portFromEnv != "" {
	//	defaultPort, _ = strconv.Atoi(portFromEnv)
	//}
	//flag.IntVar(&serverConfig.port, "listenPort", defaultPort, "The port to listen on")
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
	//viper.BindPFlags(flag.CommandLine)
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
