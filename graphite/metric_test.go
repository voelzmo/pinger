package graphite

import (
	"testing"
	"time"
)

type event struct {
	metric string
	value float64
	when int64
}

type senderSpy struct {
  sent []event
}

func (s *senderSpy) Send(metric string, value float64, when int64) error {
	s.sent = append(s.sent, event{metric, value, when})
	return nil
}

func (s senderSpy) getSent() []event {
	return s.sent
}


func TestMetricSent(t *testing.T) {
     sender := senderSpy{}
     metric := NewMetric("my.prefix", 1, &sender)

     metric.Increment()
     metric.Increment()
     metric.Increment()

     time.Sleep(3 * time.Second)

    events := sender.getSent()
	count := len(events)
	expected := 2
	if count < expected  {
		t.Errorf("Number of received events wrong, got: %d, want at least: %d.", len(events), expected)
	}

	firstMetric := events[0]

	if firstMetric.metric != "my.prefix.pingReceiveRate"  {
		t.Errorf("Wrong metric received: %s", firstMetric.metric)
	}

	if firstMetric.value != 3.0  {
		t.Errorf("Wrong metric value received: %f", firstMetric.value)
	}

	secondMetric := events[1]

	if firstMetric.when > secondMetric.when  {
		t.Errorf("Second event should occur after first: %d, %d", firstMetric.when, secondMetric.when)
	}

	secAfter1970 := firstMetric.when
	if (secAfter1970 <  1500000000) || (secAfter1970 >  6600000000) {
		t.Errorf("Time value is suspicious: %d", secAfter1970)
	}

}