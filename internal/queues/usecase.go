package queues

import (
	"context"

	"github.com/VladSatyshev/concurrent-queue/internal/models"
)

type UseCase interface {
	GetByName(ctx context.Context, queueName string) (models.Queue, error)
	GetAll(ctx context.Context) []models.Queue
	AddMessage(ctx context.Context, queueName string, jsonBody map[string]interface{}) error
	AddSubscriber(ctx context.Context, queueName string, subscriberName string) error
	ConsumeMessages(ctx context.Context, queueName string, subscriberName string) (map[string]interface{}, error)
}
