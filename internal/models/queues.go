package models

import (
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
	"github.com/VladSatyshev/concurrent-queue/pkg/utils"
)

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

func (q *Queue) HasSubscriber(name string) bool {
	for sub := range q.Subscribers {
		if sub == name {
			return true
		}
	}
	return false
}

func (q *Queue) GetNotSeenMessages(name string) map[string]interface{} {
	res := map[string]interface{}{}

	for messageID, message := range q.Messages {
		if _, ok := message.SeenBy[name]; !ok {
			res[messageID] = message
		}
	}

	return res
}

func (q *Queue) SetMessagesSeenBy(name string) {
	for _, message := range q.Messages {
		if _, ok := message.SeenBy[name]; !ok {
			message.SeenBy[name] = struct{}{}
		}
	}
}

func (q *Queue) DeleteSeenByAllMessages(logger logger.Logger) {
	for messageID, message := range q.Messages {
		if len(message.SeenBy) == len(q.Subscribers) {
			logger.Warnf("message with message ID %s has been deleted from queue %s", messageID, q.Name)
			delete(q.Messages, messageID)
		}
	}
}
