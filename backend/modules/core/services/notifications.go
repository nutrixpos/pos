package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/google/uuid"
	"github.com/olahol/melody"
)

type INotificationService interface {
	HandleHttpRequest(w http.ResponseWriter, r *http.Request) error
	SendToTopic(topic_name string, message string) error
}

type MelodyWebsocket struct {
	Logger logger.ILogger
	Config config.Config
	melody *melody.Melody
	Topics []models.Topic
}

func (ws *MelodyWebsocket) HandleHttpRequest(w http.ResponseWriter, r *http.Request) error {

	err := ws.melody.HandleRequest(w, r)
	if err != nil {
		return err
	}

	return nil
}

func (ws *MelodyWebsocket) SendToTopic(topic_name string, message string) error {

	for _, topic := range ws.Topics {
		if topic.Name == topic_name || topic.Name == "all" {
			for _, subscriber := range topic.Subscribers {
				ws.SendToSession(message, subscriber)
			}
		}
	}

	return nil
}

func (ws *MelodyWebsocket) SendToSession(msg string, session_id string) {

	ws.melody.BroadcastFilter([]byte(msg), func(q *melody.Session) bool {

		if sessionId, exists := q.Get("sessionID"); exists {
			return sessionId.(string) == session_id
		}

		return false

	})
}

func (ws *MelodyWebsocket) HandleConnect() {
	ws.melody.HandleConnect(func(s *melody.Session) {
		sessionID := uuid.New().String()
		s.Set("sessionID", sessionID) // Store the ID in the session context
	})
}

func (ws *MelodyWebsocket) AddSessionToTopic(topic_name string, session_id string) {
	if _, index, err := ws.GetTopic(topic_name); err == nil {
		ws.Topics[index].Subscribers = append(ws.Topics[index].Subscribers, session_id)
	} else if err.Error() == "topic not found" {
		ws.Topics = append(ws.Topics, models.Topic{
			Name:        topic_name,
			Subscribers: []string{session_id},
		})
	}
}

func (ws *MelodyWebsocket) HandleMessages() {
	ws.melody.HandleMessage(func(s *melody.Session, msg []byte) {

		session_id, exists := s.Get("sessionID")
		if !exists {
			ws.SendToSession("{state:\"connection not found\"}", session_id.(string))
			return
		}

		var message models.WebsocketClientBaseMessage
		if err := json.Unmarshal(msg, &message); err != nil {
			ws.Logger.Error(err.Error())
			return
		}

		if message.Type == "subscribe" {

			var subscribe_message models.WebsocketSubscribeClientMessage
			if err := json.Unmarshal([]byte(msg), &subscribe_message); err != nil {
				ws.Logger.Error(err.Error())
				return
			}

			ws.AddSessionToTopic(subscribe_message.TopicName, session_id.(string))

		}

		if message.Type == "topic_message" {

			var topic_message models.WebsocketTopicClientMessage
			if err := json.Unmarshal(msg, &topic_message); err != nil {
				ws.Logger.Error(err.Error())
				return
			}

			if topic_message.TopicName == "order_finished" {

				var order_finish_client_message models.WebsocketOrderFinishClientMessage
				if err := json.Unmarshal([]byte(msg), &order_finish_client_message); err != nil {
					ws.Logger.Error(err.Error())
					return
				}

				order_finish_topic_message := models.WebsocketOrderFinishServerMessage{
					WebsocketTopicServerMessage: models.WebsocketTopicServerMessage{
						Type:      "topic_message",
						TopicName: "order_finished",
						Severity:  "info",
					},
					OrderId: order_finish_client_message.OrderId,
				}

				order_finish_topic_message_json, err := json.Marshal(order_finish_topic_message)
				if err != nil {
					ws.Logger.Error(err.Error())
					return
				}

				ws.SendToTopic("order_finish", string(order_finish_topic_message_json))
				ws.SendToSession("{state:\"success\"}", session_id.(string))
			}
		}

	})
}

func (ws *MelodyWebsocket) GetTopic(topic_name string) (topic models.Topic, index int, err error) {

	for _, t := range ws.Topics {
		if t.Name == topic_name {
			return t, 0, nil
		}
	}

	return topic, 0, fmt.Errorf("topic not found")
}

var melody_ws *MelodyWebsocket
var once sync.Once

func SpawnNotificationService(name string, logger logger.ILogger, config config.Config) (INotificationService, error) {

	switch name {
	case "melody":

		once.Do(func() {
			melody_ws = &MelodyWebsocket{
				Logger: logger,
				Config: config,
				melody: melody.New(),
			}

			melody_ws.HandleConnect()
			melody_ws.HandleMessages()

		})

		return melody_ws, nil

	default:
		return nil, fmt.Errorf("unknown notification service name: %s", name)
	}
}
