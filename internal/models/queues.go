package models

import "github.com/VladSatyshev/concurrent-queue/pkg/utils"

type Queue struct {
	Name           string
	MaxLength      uint
	MaxSubscribers uint
	Subscribers    map[string]struct{}
	Messages       map[string]QueueMessage
}

type QueueMessage struct {
	Body   map[string]interface{}
	SeenBy map[string]struct{}
}

func (q *Queue) AddMessage(jsonBody map[string]interface{}) {
	messageID := utils.GenerateMessageID()

	q.Messages[messageID] = QueueMessage{
		Body:   jsonBody,
		SeenBy: map[string]struct{}{},
	}
}

func (q *Queue) AddSubscriber(name string) {
	q.Subscribers[name] = struct{}{}
}
