package kafka

const (
	subscriptionTopic = "subscriptions"
	subscriptionType  = "subscription"
	postTopic         = "posts"
	postType          = "post"
)

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
