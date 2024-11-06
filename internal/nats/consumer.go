package nats

import "github.com/nats-io/nats.go"

type Consumer struct {
	conn *nats.Conn
}

func NewConsumer(conn *nats.Conn) Consumer {
	return Consumer{conn: conn}
}

type HandlerFunc func(data []byte)

func (c *Consumer) Subscribe(subj string, h HandlerFunc) error {
	f := func(msg *nats.Msg) { h(msg.Data) }
	_, err := c.conn.Subscribe(subj, f)
	return err
}
