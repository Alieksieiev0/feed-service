package kafka

import (
	"context"
	"encoding/json"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(receivers []models.User, notif Notification) error
}

type producer struct {
	addr string
}

func NewProducer(addr string) Producer {
	return &producer{
		addr: addr,
	}
}

func (p *producer) Produce(receivers []models.User, notif Notification) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.addr),
		Topic:                  notif.Topic(),
		AllowAutoTopicCreation: true,
	}

	value, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	messages := []kafka.Message{}
	for _, r := range receivers {
		message := kafka.Message{
			Key:   []byte(r.ID),
			Value: value,
		}

		messages = append(messages, message)
	}

	if err = w.WriteMessages(context.Background(), messages...); err != nil {
		return err
	}

	return w.Close()
}
