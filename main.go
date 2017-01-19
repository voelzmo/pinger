package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type addressVar []string

var (
	interval        int
	addressesToPing addressVar
	listenPort      int
	errorRate       float64
	randomGenerator *rand.Rand
)

func (a *addressVar) String() string {
	return fmt.Sprint(*a)
}

func (a *addressVar) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func pingAddress(address string) {
	response, err := http.Get("http://" + address + "/ping")

	if err != nil {
		log.Printf("Couldn't ping %s: %s", address, err)
	} else {
		defer response.Body.Close()
		var pingAnswer string
		json.NewDecoder(response.Body).Decode(&pingAnswer)
		log.Printf("Pinged %s, got answer: %s", address, pingAnswer)
	}
}

func startHTTPServer() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		f := randomGenerator.Float64()
		if f < errorRate {
			http.Error(w, "no wai!", http.StatusTeapot)
			return
		}
		json.NewEncoder(w).Encode("pong")
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", listenPort), nil))
}

func init() {
	flag.IntVar(&interval, "interval", 10, "The interval between two pings in seconds")
	flag.Var(&addressesToPing, "address", "The address which should be pinged. Format <IP>:<port>")
	flag.Float64Var(&errorRate, "error-rate", 0.0, "error rate in percent")

	portFromEnv := os.Getenv("PORT")
	defaultPort := 8080
	if portFromEnv != "" {
		defaultPort, _ = strconv.Atoi(portFromEnv)
	}
	flag.IntVar(&listenPort, "listenPort", defaultPort, "The port to listen on")
}

func main() {
	flag.Parse()
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

	go startHTTPServer()
	c := time.Tick(time.Duration(interval) * time.Second)
	for range c {
		for _, address := range addressesToPing {
			pingAddress(address)
		}
	}
}
