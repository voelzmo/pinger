package main

import (
	goflag "flag"
	"io/ioutil"
	"log"

	"github.com/facebookgo/pidfile"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configPath string
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

func init() {
	flag.StringVar(&configPath, "config-path", "", "Absolute path to read a config file from")
	flag.StringSlice("address", []string{"http://localhost:8080"}, "Server address and port to ping from the client")
	flag.Int("port", 8080, "Server port to listen on")
	flag.Int("interval", 10, "Interval in seconds for clients calling '/ping'")
	flag.Float64("error-rate", 0.0, "Server error-rate in percent. Server will respond with an error instead of a regular 'pong' answer")
}

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Fatalf("fatal error during binding flags: '%v'", err)
	}
	viper.AutomaticEnv()

	readConfiguration()
	writePIDFile(viper.GetString("pidfile"))

	serverConfig := &HTTPConfig{
		viper.GetInt("port"),
		viper.GetString("cert-path"),
		viper.GetString("key-path"),
		viper.GetString("ca-cert-path"),
	}
	verifyHTTPSConfig(serverConfig)

	pingServer := NewPingServer(serverConfig, viper.GetFloat64("error-rate"), viper.GetString("CF_INSTANCE_GUID"))
	go pingServer.Start()

	pingClient := NewPingClient(serverConfig, viper.GetStringSlice("address"), viper.GetDuration("interval"))
	go pingClient.Start()

	select {} // keep main running, as we only have gofuncs above
}

func readConfiguration() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("fatal error config file at '%s': %s", configPath, err)
		}
		content, _ := ioutil.ReadFile(configPath)
		log.Printf("config file:'%s'\n", content)
	}
}

func writePIDFile(path string) {
	if path != "" {
		pidfile.SetPidfilePath(path)
		err := pidfile.Write()
		if err != nil {
			log.Fatalf("couldn't write Pidfile at path '%s': %s", path, err)
		}
	}
}

func verifyHTTPSConfig(sc *HTTPConfig) {
	if sc.certFile != "" && sc.keyFile == "" {
		log.Fatal("specify either certificate and key or none")
	}
	if sc.certFile == "" && sc.keyFile != "" {
		log.Fatal("specify either certificate and key or none")
	}
}
