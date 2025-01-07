package usecase

import (
	"context"
	"fmt"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/models"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/VladSatyshev/concurrent-queue/internal/queues/mock"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
	"github.com/golang/mock/gomock"
)

func configureEnvironment(ctrl *gomock.Controller, queuesCfg []config.QueueConfig) (queues.UseCase, func()) {
	cfg := &config.Config{
		Logger: config.LoggerConfig{Development: true, DisableCaller: false, DisableStacktrace: false, Encoding: "json"},
		Queues: queuesCfg,
	}
	apiLogger := logger.NewAPILogger(cfg)
	apiLogger.InitLogger()
	mockQueueRepo := mock.NewMockRepository(ctrl)

	mockQueuesStorage := make(map[string]models.Queue)

	mockQueueRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, name string, maxLength uint, maxSubscribers uint) (models.Queue, error) {
		if _, ok := mockQueuesStorage[name]; ok {
			return models.Queue{}, queues.NewQueueErr(queues.RepositoryErr, fmt.Sprintf("queue with name %v already exists in mockQueuesStorage", name))
		}

		newQueue := models.Queue{
			Name:           name,
			MaxLength:      maxLength,
			MaxSubscribers: maxSubscribers,
			Subscribers:    make(map[string]struct{}, maxSubscribers),
			Messages:       make(map[string]models.QueueMessage, maxLength),
		}

		mockQueuesStorage[name] = newQueue
		return newQueue, nil
	})

	mockQueueRepo.EXPECT().GetByName(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, name string) (models.Queue, error) {
		queue, ok := mockQueuesStorage[name]
		if !ok {
			return models.Queue{}, queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", queue.Name))
		}

		return queue, nil
	})

	mockQueueRepo.EXPECT().GetAll(gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context) []models.Queue {
		res := make([]models.Queue, 0, len(mockQueuesStorage))

		for _, queue := range mockQueuesStorage {
			res = append(res, queue)
		}

		return res
	})

	mockQueueRepo.EXPECT().AddMessage(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, name string, jsonBody map[string]interface{}) error {
		q, ok := mockQueuesStorage[name]
		if !ok {
			return queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", name))
		}

		q.AddMessage(jsonBody)

		return nil
	})

	mockQueueRepo.EXPECT().AddSubscriber(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, queueName string, subscriberName string) error {
		q, ok := mockQueuesStorage[queueName]
		if !ok {
			return queues.NewQueueErr(queues.RepositoryNotFoundErr, fmt.Sprintf("queue %s not found", queueName))
		}

		q.Subscribers[subscriberName] = struct{}{}

		return nil
	})

	for _, q := range queuesCfg {
		_, err := mockQueueRepo.Create(context.Background(), q.Name, q.Length, q.SubscribersAmount)
		if err != nil {
			panic(err)
		}
	}

	queuesUC := NewQueuesUseCase(cfg, mockQueueRepo, apiLogger)

	return queuesUC, func() {
		// cleanup calls go here
		ctrl.Finish()
	}
}
