package graphite

import (
	"fmt"
	"net"
	"time"
)

const (
	tcpTimeout = 30 * time.Second
)

type GraphiteSender struct {
	connection net.Conn
}

func NewGraphiteSender(endpoint string) (*GraphiteSender, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &GraphiteSender{connection: conn}, nil
}

func (gs *GraphiteSender) Send(metric string, value float64, when int64) error {
	err := gs.connection.SetWriteDeadline(time.Now().Add(tcpTimeout))
	if err != nil {
		return err
	}
	_, err = gs.connection.Write([]byte(fmt.Sprintf("%s %f %d\n", metric, value, when)))
	return err
}

func (gs *GraphiteSender) Close() error {
	return gs.connection.Close()
}
