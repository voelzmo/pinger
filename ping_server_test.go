package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/voelzmo/pinger"
)

var _ = Describe("PingServer", func() {
	var pingServer = NewPingServer(&HTTPConfig{}, 0.0)

	Describe("Can be started", func() {
		pingServer.Start()
	})
})
