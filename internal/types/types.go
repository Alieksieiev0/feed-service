package types

import (
	"encoding/xml"
	"time"
)

type SubscriptionPartialSuccess struct {
	XMLName      xml.Name    `xml:"multistatus"`
	Subscription XMLResponse `xml:"subscription"`
	Notification XMLResponse `xml:"notification"`
}

type PostPartialSuccess struct {
	XMLName      xml.Name    `xml:"multistatus"`
	Creation     XMLResponse `xml:"creation"`
	Notification XMLResponse `xml:"notification"`
}

type XMLResponse struct {
	Status int    `xml:"status"`
	Error  string `xml:"error"`
}

type Post struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	OwnerName string    `json:"owner_name"`
	OwnerId   string    `json:"owner_id"`
}

type UserBase struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	UserBase
	Password    string     `json:"password"`
	Subscribers []UserBase `json:"subcribers"`
	Posts       []Post     `json:"posts"`
}
