package graphite

import (
	"time"
	"sync/atomic"
)

type Metric struct {
	pingEvents  int32
	metricPrefix string
	sendInterval float64
	sender       Sender
}

func NewMetric(metricPrefix string, sendInterval float64, sender Sender) *Metric {
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
	ticker := time.Tick(time.Duration(m.sendInterval) * time.Second)
	for now := range ticker {
		for {
			currentValue := m.pingEvents
			if atomic.CompareAndSwapInt32(&m.pingEvents, currentValue, 0) {
				m.sender.Send(m.metricPrefix + ".pingReceiveRate", float64(currentValue)/m.sendInterval, now.Unix())
				break
			}
		}
	}
}
