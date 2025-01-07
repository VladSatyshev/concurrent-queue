package queues

import (
	"context"

	"github.com/VladSatyshev/concurrent-queue/internal/models"
)

//go:generate mockgen -source repository.go -destination mock/repository_mock.go -package mock
type Repository interface {
	Create(ctx context.Context, name string, maxLength uint, maxSubscribers uint) (models.Queue, error)
	GetByName(ctx context.Context, name string) (models.Queue, error)
	GetAll(ctx context.Context) []models.Queue
	AddMessage(ctx context.Context, name string, jsonBody map[string]interface{}) error
	AddSubscriber(ctx context.Context, queueName string, subscriberName string) error
}
