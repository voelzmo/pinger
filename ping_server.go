package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// PingServer is a webserver that can be started.
type PingServer interface {
	Start()
}

type pingServer struct {
	config          *HTTPConfig
	randomGenerator *rand.Rand
	errorRate       float64
}

// NewPingServer creates a new PingServer obviously.
func NewPingServer(serverConfig *HTTPConfig, errorRate float64) PingServer {
	return &pingServer{
		config:          serverConfig,
		randomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
		errorRate:       errorRate,
	}
}

func (ps *pingServer) Start() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		f := ps.randomGenerator.Float64()
		if f < errorRate {
			http.Error(w, "no wai!", http.StatusTeapot)
			return
		}
		json.NewEncoder(w).Encode("pong")
		log.Printf("Got pinged by %s, sent a 'pong'!", r.RemoteAddr)
	})
	log.Fatal(ps.listenAndServe())
}

func (ps *pingServer) listenAndServe() error {
	address := fmt.Sprintf(":%v", ps.config.port)
	if !ps.config.isSecure() {
		return http.ListenAndServe(address, nil)
	}
	return http.ListenAndServeTLS(address, ps.config.certFile, ps.config.keyFile, nil)
}
