package models

type Topic struct {
	Name        string
	Subscribers []string
}

type WebsocketClientBaseMessage struct {
	Type string `json:"type"`
}

type WebsocketSubscribeClientMessage struct {
	WebsocketClientBaseMessage `json:",inline"`
	TopicName                  string `json:"topic_name"`
}

type WebsocketTopicClientMessage struct {
	TopicName string `json:"topic_name"`
}

type WebsocketOrderFinishClientMessage struct {
	OrderId string `json:"order_id"`
}

type WebsocketTopicServerMessage struct {
	Type      string `json:"type"`
	TopicName string `json:"topic_name"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Key       string `json:"key"` // unqique label to prevent message duplications
}

type WebsocketOrderFinishServerMessage struct {
	WebsocketTopicServerMessage `json:",inline"`
	OrderId                     string `json:"order_id"`
}
