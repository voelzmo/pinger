package graphite

import (
	"fmt"
	"net"
	"time"
)

const (
	tcpTimeout = 30 * time.Second
)

type Sender interface {
	Send(metric string, value float64, when int64) error
}

type graphiteSender struct {
	endpoint string
}

func NewGraphiteSender(endpoint string) Sender {
	return &graphiteSender{endpoint}
}

func (gs *graphiteSender) Send(metric string, value float64, when int64) error {
	conn, err := net.Dial("tcp", gs.endpoint)
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
	}()

	err = conn.SetWriteDeadline(time.Now().Add(tcpTimeout))
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%s %f %d\n", metric, value, when)))
	return err
}
