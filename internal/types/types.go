package types

import "encoding/xml"

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
