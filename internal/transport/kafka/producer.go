package kafka

import (
	"context"
	"encoding/json"

	"github.com/Alieksieiev0/feed-service/internal/types"
	"github.com/segmentio/kafka-go"
)

const (
	subscriptionTopic = "subscriptions"
	subscriptionType  = "subscription"
	postTopic         = "posts"
	postType          = "post"
)

type Producer interface {
	Produce(receivers []types.UserBase, notification Notification) error
}

type producer struct {
	addr string
}

func NewProducer(addr string) Producer {
	return &producer{
		addr: addr,
	}
}

func (p *producer) Produce(receivers []types.UserBase, notification Notification) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.addr),
		Topic:                  notification.Topic(),
		AllowAutoTopicCreation: true,
	}

	value, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	messages := []kafka.Message{}
	for _, r := range receivers {
		message := kafka.Message{
			Key:   []byte(r.Id),
			Value: value,
		}

		messages = append(messages, message)
	}

	if err = w.WriteMessages(context.Background(), messages...); err != nil {
		return err
	}

	return w.Close()
}

type Notification interface {
	Topic() string
}

type notification struct {
	FromId   string
	FromName string
	TargetId string
	Type     string
	topic    string
}

func (n *notification) Topic() string {
	return n.topic
}

func NewPostNotification(ownerId, ownerName, postId string) Notification {
	return &notification{
		FromId:   ownerId,
		FromName: ownerName,
		TargetId: postId,
		Type:     postType,
		topic:    postTopic,
	}
}

func NewSubscriptionNotification(userId, userName string) Notification {
	return &notification{
		FromId:   userId,
		FromName: userName,
		Type:     subscriptionType,
		topic:    postTopic,
	}
}
