package kafka

import (
	"context"
	"encoding/json"
	"fmt"

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
	Produce(receivers []types.UserBase, message Message) error
}

type producer struct {
	addr string
}

func NewProducer(addr string) Producer {
	return &producer{
		addr: addr,
	}
}

func (p *producer) Produce(receivers []types.UserBase, message Message) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(p.addr),
		Topic:                  message.Topic(),
		AllowAutoTopicCreation: true,
	}

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}
	fmt.Println("---")
	fmt.Printf("receiver: %+v\n ", receivers[0])
	fmt.Printf("msg: %+v \n", message)

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

type Message interface {
	Topic() string
}

type message struct {
	FromId   string `json:"from_id"`
	FromName string `json:"from_name"`
	TargetId string `json:"target_id"`
	Type     string `json:"type"`
	topic    string
}

func (n *message) Topic() string {
	return n.topic
}

func NewPostMessage(ownerId, ownerName, postId string) Message {
	return &message{
		FromId:   ownerId,
		FromName: ownerName,
		TargetId: postId,
		Type:     postType,
		topic:    postTopic,
	}
}

func NewSubscriptionMessage(userId, userName string) Message {
	return &message{
		FromId:   userId,
		FromName: userName,
		Type:     subscriptionType,
		topic:    postTopic,
	}
}
