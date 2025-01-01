package usecase

import (
	"context"
	"fmt"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/models"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
)

type queuesUC struct {
	cfg        *config.Config
	queuesRepo queues.Repository
	logger     logger.Logger
}

func NewQueuesUseCase(cfg *config.Config, queuesRepo queues.Repository, logger logger.Logger) queues.UseCase {
	return &queuesUC{
		cfg:        cfg,
		queuesRepo: queuesRepo,
		logger:     logger,
	}
}

// get queue by name
func (u *queuesUC) GetByName(ctx context.Context, name string) (models.Queue, error) {
	queue, err := u.queuesRepo.GetByName(ctx, name)
	if err != nil {
		return models.Queue{}, err
	}

	return queue, nil
}

// get all queues
func (u *queuesUC) GetAll(ctx context.Context) ([]models.Queue, error) {
	queues, err := u.queuesRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return queues, nil
}

func (u *queuesUC) AddMessage(ctx context.Context, name string, jsonBody map[string]interface{}) error {
	queue, err := u.queuesRepo.GetByName(ctx, name)
	if err != nil {
		return err
	}

	if len(queue.Messages) == int(queue.MaxLength) {
		return fmt.Errorf("too many messages: max amount of messages for queue %v is %v", name, queue.MaxLength)
	}

	queue.AddMessage(jsonBody)

	return nil
}

func (u *queuesUC) AddSubscriber(ctx context.Context, queueName string, subscriberName string) error {
	queue, err := u.queuesRepo.GetByName(ctx, queueName)
	if err != nil {
		return err
	}

	if len(queue.Subscribers) == int(queue.MaxSubscribers) {
		return fmt.Errorf("too many subscribers: max amount of subscribers for queue %v is %v", queueName, queue.MaxSubscribers)
	}

	queue.AddSubscriber(subscriberName)

	return nil
}
