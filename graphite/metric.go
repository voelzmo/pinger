package graphite

import (
	"log"
	"sync/atomic"
	"time"
)

type Metric struct {
	pingEvents   int32
	metricPrefix string
	sendInterval time.Duration
	sender       Sender
}

func NewMetric(metricPrefix string, sendInterval time.Duration, sender Sender) *Metric {
	result := Metric{
		metricPrefix: metricPrefix,
		sendInterval: sendInterval,
		sender:       sender,
	}
	go result.reportMetrics()
	return &result
}

func (m *Metric) Increment() {
	atomic.AddInt32(&m.pingEvents, 1)
}

func (m *Metric) reportMetrics() {
	ticker := time.Tick(m.sendInterval)
	for now := range ticker {
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
