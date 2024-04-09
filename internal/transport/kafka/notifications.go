package kafka

const (
	subscriptionTopic = "subscriptions"
	postTopic         = "posts"
)

type Notification interface {
	Topic() string
}

type Base struct {
	FromId   string
	FromName string
	topic    string
}

func (bn *Base) Topic() string {
	return bn.topic
}

func NewPost(fromId, fromName, postId string, topic ...string) *Post {
	t := postTopic
	if len(topic) > 0 {
		t = topic[0]
	}

	return &Post{
		Base: Base{
			FromId:   fromId,
			FromName: fromName,
			topic:    t,
		},
		PostId: postId,
	}
}

type Post struct {
	Base
	PostId string
}

func NewSubscription(fromId, fromName string, topic ...string) *Subscription {
	t := subscriptionTopic
	if len(topic) > 0 {
		t = topic[0]
	}

	return &Subscription{
		Base: Base{
			FromId:   fromId,
			FromName: fromName,
			topic:    t,
		},
	}
}

type Subscription struct {
	Base
}
