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

func NewQueuesRepository(cfg *config.Config) (queues.Repository, error) {
	resQueues := make(map[string]models.Queue, len(cfg.Queues))

	for _, queue := range cfg.Queues {
		if _, ok := resQueues[queue.Name]; ok {
			return nil, queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue with name %v already exist", queue.Name))
		}

		newQueue := models.Queue{
			Name:           queue.Name,
			MaxLength:      queue.Length,
			MaxSubscribers: queue.SubscribersAmount,
			Subscribers:    make(map[string]struct{}, queue.SubscribersAmount),
			Messages:       make(map[string]models.QueueMessage, queue.Length),
		}

		resQueues[queue.Name] = newQueue
	}

	return &queuesRepo{queues: resQueues}, nil
}

func (r *queuesRepo) GetByName(ctx context.Context, name string) (models.Queue, error) {
	queue, ok := r.queues[name]
	if !ok {
		return models.Queue{}, queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue %s not found", queue.Name))
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
		return queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue %s not found", name))
	}

	q.AddMessage(jsonMsgBody)

	return nil
}

func (r *queuesRepo) AddSubscriber(ctx context.Context, queueName string, subscriberName string) error {
	q, ok := r.queues[queueName]
	if !ok {
		return queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue %s not found", queueName))
	}

	q.Subscribers[subscriberName] = struct{}{}

	return nil
}
