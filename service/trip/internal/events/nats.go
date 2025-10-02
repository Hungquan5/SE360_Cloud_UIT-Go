package events

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	SubTripRequested  = "trip.requested"
	SubDriverAccepted = "driver.accepted"
	SubDriverLocation = "driver.location.updated"
	SubTripCompleted  = "trip.completed" // optional publish on complete
)

type Bus struct {
	NATS *nats.Conn
}

func Connect(url string) (*Bus, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Bus{NATS: nc}, nil
}

func (b *Bus) PublishJSON(subject string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return b.NATS.Publish(subject, data)
}

func (b *Bus) Subscribe(subject string, fn func(raw []byte)) (*nats.Subscription, error) {
	return b.NATS.Subscribe(subject, func(m *nats.Msg) {
		fn(m.Data)
	})
}

func LogErr(err error) {
	if err != nil {
		log.Println("nats:", err)
	}
}
