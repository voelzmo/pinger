package graphite

import (
	"code.cloudfoundry.org/clock"
	"log"
	"sync/atomic"
	"time"
)

type Metric struct {
	pingEvents   int32
	metricPrefix string
	sendInterval time.Duration
	sender       Sender
	c            clock.Clock
}

func NewMetric(metricPrefix string, sendInterval time.Duration, sender Sender, c clock.Clock) *Metric {
	result := Metric{
		metricPrefix: metricPrefix,
		sendInterval: sendInterval,
		sender:       sender,
		c:            c,
	}
	go result.reportMetrics()
	return &result
}

func (m *Metric) Increment() {
	atomic.AddInt32(&m.pingEvents, 1)
}

func (m *Metric) reportMetrics() {
	ticker := m.c.NewTicker(m.sendInterval)
	for now := range ticker.C() {
		for {
			currentValue := m.pingEvents
			if atomic.CompareAndSwapInt32(&m.pingEvents, currentValue, 0) {
				rate := float64(currentValue) / m.sendInterval.Seconds()
				err := m.sender.Send(m.metricPrefix+".pingReceiveRate", rate, now.Unix())
				if nil != err {
					log.Printf("Error while sending metric: %s", err.Error())
				}
				break
			}
		}
	}
}
