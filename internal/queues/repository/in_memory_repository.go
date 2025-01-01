package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/models"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
)

type queuesRepo struct {
	queues map[string]models.Queue
}

func NewQueuesRepository(cfg *config.Config) (queues.Repository, error) {
	queues := make(map[string]models.Queue, len(cfg.Queues))

	for _, queue := range cfg.Queues {
		if _, ok := queues[queue.Name]; ok {
			return nil, fmt.Errorf("queue with name %v already exist", queue.Name)
		}

		newQueue := models.Queue{
			Name:           queue.Name,
			MaxLength:      queue.Length,
			MaxSubscribers: queue.SubscribersAmount,
			Subscribers:    make(map[string]struct{}, queue.SubscribersAmount),
			Messages:       make(map[string]models.QueueMessage, queue.Length),
		}

		queues[queue.Name] = newQueue
	}

	return &queuesRepo{queues: queues}, nil
}

func (r *queuesRepo) GetByName(ctx context.Context, name string) (models.Queue, error) {
	queue, ok := r.queues[name]
	if !ok {
		return models.Queue{}, errors.New("queue not found")
	}

	return queue, nil
}

func (r *queuesRepo) GetAll(ctx context.Context) ([]models.Queue, error) {
	res := make([]models.Queue, 0, len(r.queues))

	for _, queue := range r.queues {
		res = append(res, queue)
	}

	return res, nil
}

func (r *queuesRepo) AddMessage(ctx context.Context, name string, jsonMsgBody map[string]interface{}) error {
	q, ok := r.queues[name]
	if !ok {
		return errors.New("queue not found")
	}

	q.AddMessage(jsonMsgBody)

	return nil
}

func (r *queuesRepo) AddSubscriber(ctx context.Context, queueName string, subscriberName string) error {
	q, ok := r.queues[queueName]
	if !ok {
		return errors.New("queue not found")
	}

	q.Subscribers[subscriberName] = struct{}{}

	return nil
}
