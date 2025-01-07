package repository

import (
	"context"
	"fmt"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/models"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
)

type queuesRepo struct {
	queues map[string]models.Queue
}

func NewQueuesRepository(cfg *config.Config) queues.Repository {
	resQueues := make(map[string]models.Queue, len(cfg.Queues))
	return &queuesRepo{queues: resQueues}
}

func InitQueues(ctx context.Context, cfg *config.Config, r queues.Repository) error {
	for _, queue := range cfg.Queues {
		if _, err := r.Create(ctx, queue.Name, queue.Length, queue.SubscribersAmount); err != nil {
			return err
		}
	}

	return nil
}

func (r *queuesRepo) Create(ctx context.Context, name string, maxLength uint, maxSubscribers uint) (models.Queue, error) {
	if _, ok := r.queues[name]; ok {
		return models.Queue{}, queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue with name %v already exists", name))
	}

	newQueue := models.Queue{
		Name:           name,
		MaxLength:      maxLength,
		MaxSubscribers: maxSubscribers,
		Subscribers:    make(map[string]struct{}, maxSubscribers),
		Messages:       make(map[string]models.QueueMessage, maxLength),
	}

	r.queues[name] = newQueue

	return newQueue, nil
}

func (r *queuesRepo) GetByName(ctx context.Context, name string) (models.Queue, error) {
	queue, ok := r.queues[name]
	if !ok {
		return models.Queue{}, queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", queue.Name))
	}

	return queue, nil
}

func (r *queuesRepo) GetAll(ctx context.Context) []models.Queue {
	res := make([]models.Queue, 0, len(r.queues))

	for _, queue := range r.queues {
		res = append(res, queue)
	}

	return res
}

func (r *queuesRepo) AddMessage(ctx context.Context, name string, jsonMsgBody map[string]interface{}) error {
	q, ok := r.queues[name]
	if !ok {
		return queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", name))
	}

	q.AddMessage(jsonMsgBody)

	return nil
}

func (r *queuesRepo) AddSubscriber(ctx context.Context, queueName string, subscriberName string) error {
	q, ok := r.queues[queueName]
	if !ok {
		return queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", queueName))
	}

	q.Subscribers[subscriberName] = struct{}{}

	return nil
}
