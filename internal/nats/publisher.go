package nats

import (
	"encoding/json"
	"wb-challenge/internal"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) Publisher {
	return Publisher{conn: conn}
}

func (p *Publisher) Publish(events ...internal.Event) error {
	for _, e := range events {
		raw, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if err := p.conn.Publish(e.Type(), raw); err != nil {
			return err
		}
	}
	return nil
}
