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

func (u *queuesUC) getByName(ctx context.Context, name string) (models.Queue, error) {
	queue, err := u.queuesRepo.GetByName(ctx, name)
	if err != nil {
		if qErr, ok := err.(*queues.QueueErr); ok {
			if qErr.ErrType == queues.RepositoryNotFoundErr {
				u.logger.Errorf("queue %s was not found", name)
				return models.Queue{}, queues.NewQueueErr(queues.UseCaseNotFoundErr, "Queue not found")
			}
		}
		return models.Queue{}, err
	}
	return queue, nil
}

// get queue by name
func (u *queuesUC) GetByName(ctx context.Context, name string) (models.Queue, error) {
	u.logger.Info("GetByName UC is in action")
	queue, err := u.getByName(ctx, name)
	if err != nil {
		return models.Queue{}, err
	}

	return queue, nil
}

// get all queues
func (u *queuesUC) GetAll(ctx context.Context) []models.Queue {
	u.logger.Info("GetAll UC is in action")
	queues := u.queuesRepo.GetAll(ctx)
	return queues
}

// add message to queue
func (u *queuesUC) AddMessage(ctx context.Context, name string, jsonBody map[string]interface{}) error {
	u.logger.Info("AddMessage UC is in action")
	queue, err := u.getByName(ctx, name)
	if err != nil {
		return err
	}

	if len(queue.Messages) >= int(queue.MaxLength) {
		msg := "too many messages: max amount of messages for queue %v is %v"
		u.logger.Errorf(msg, name, queue.MaxLength)
		return queues.NewQueueErr(queues.UseCaseErr, fmt.Sprintf(msg, name, queue.MaxLength))
	}

	queue.AddMessage(jsonBody)

	u.logger.Infof("Message %v has been added to queue %s", jsonBody, queue.Name)

	return nil
}

// add subscriber to queue
func (u *queuesUC) AddSubscriber(ctx context.Context, queueName string, subscriberName string) error {
	u.logger.Info("AddSubscriber UC is in action")
	queue, err := u.getByName(ctx, queueName)
	if err != nil {
		return err
	}

	if queue.HasSubscriber(subscriberName) {
		return queues.NewQueueErr(queues.UseCaseErr, fmt.Sprintf("user %s has already subscribed to queue %s", subscriberName, queue.Name))
	}

	if len(queue.Subscribers) == int(queue.MaxSubscribers) {
		return queues.NewQueueErr(queues.UseCaseErr, fmt.Sprintf("too many subscribers: max amount of subscribers for queue %v is %v", queueName, queue.MaxSubscribers))
	}

	queue.AddSubscriber(subscriberName)

	u.logger.Infof("Subscriber %s has been added to queue %s", subscriberName, queue.Name)

	return nil
}

// consume messages from queue by subscriber
func (u *queuesUC) ConsumeMessages(ctx context.Context, queueName string, subscriberName string) (map[string]interface{}, error) {
	u.logger.Info("ConsumeMessages UC is in action")
	queue, err := u.getByName(ctx, queueName)
	if err != nil {
		return nil, err
	}

	if !queue.HasSubscriber(subscriberName) {
		return nil, queues.NewQueueErr(queues.UseCaseErr, fmt.Sprintf("queue %v doesn't have subscriber %s", queue.Name, subscriberName))
	}

	notSeenMessages := queue.GetNotSeenMessages(subscriberName)

	queue.SetMessagesSeenBy(subscriberName)

	queue.DeleteSeenByAllMessages(u.logger)

	return notSeenMessages, nil
}
