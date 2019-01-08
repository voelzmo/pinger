package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPinger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pinger Suite")
}
