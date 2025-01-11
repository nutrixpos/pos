// Package models contains the data models for the application.
//
// The models are used to store data in the database and to marshal/unmarshal
// data to/from JSON.
package models

import "time"

// Topic represents a group of subscribers for a specific topic.
type Topic struct {
	Name        string
	Subscribers []string
}

// WebsocketClientBaseMessage is a base message for all client messages.
type WebsocketClientBaseMessage struct {
	Type string `json:"type"`
}

// WebsocketSubscribeClientMessage is a message sent by the client to subscribe
// to a specific topic.
type WebsocketSubscribeClientMessage struct {
	WebsocketClientBaseMessage `json:",inline"`
	TopicName                  string `json:"topic_name"`
}

// WebsocketTopicClientMessage is a message sent by the client to send a message
// to a specific topic.
type WebsocketTopicClientMessage struct {
	TopicName string `json:"topic_name"`
}

// WebsocketOrderFinishClientMessage is a message sent by the client to finish
// an order.
type WebsocketOrderFinishClientMessage struct {
	OrderId string `json:"order_id"`
}

// WebsocketTopicServerMessage is a message sent by the server to a specific topic.
type WebsocketTopicServerMessage struct {
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	TopicName string    `json:"topic_name"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	Key       string    `json:"key"` // unqique label to prevent message duplications
}

// WebsocketOrderFinishServerMessage is a message sent by the server to finish
// an order.
type WebsocketOrderFinishServerMessage struct {
	WebsocketTopicServerMessage `json:",inline"`
	OrderId                     string `json:"order_id"`
}

type WebsocketOrderSubmitServerMessage struct {
	WebsocketTopicServerMessage `json:",inline"`
	Order                       Order `json:"order"`
}
