package usecase

import (
	"context"
	"testing"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQueuesUC_AddMessageToQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	expectedQueue := models.Queue{
		Name:           qConfig.Name,
		MaxLength:      qConfig.Length,
		MaxSubscribers: qConfig.SubscribersAmount,
	}
	msgBody := map[string]interface{}{"msg": "hello"}
	ctx := context.Background()

	err := queuesUC.AddMessage(ctx, qConfig.Name, msgBody)
	assert.Nil(t, err)

	actualQ, err := queuesUC.GetByName(ctx, qConfig.Name)
	assert.Nil(t, err)

	assert.Equal(t, expectedQueue.Name, actualQ.Name)
	assert.Equal(t, expectedQueue.MaxLength, actualQ.MaxLength)
	assert.Equal(t, expectedQueue.MaxSubscribers, actualQ.MaxSubscribers)

	assert.Equal(t, 1, len(actualQ.Messages))
	for _, m := range actualQ.Messages {
		assert.Equal(t, m.Body, msgBody)
	}
}

func TestQueuesUC_CantAddMessageToFullQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	msgBody1 := map[string]interface{}{"msg": "hello1"}
	msgBody2 := map[string]interface{}{"msg": "hello2"}
	ctx := context.Background()

	err := queuesUC.AddMessage(ctx, qConfig.Name, msgBody1)
	assert.Nil(t, err)
	err = queuesUC.AddMessage(ctx, qConfig.Name, msgBody2)
	assert.NotNil(t, err)
}

func TestQueuesUC_CanSubscribeToQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName := "subscriber"

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName)
	assert.Nil(t, err)
	q, err := queuesUC.GetByName(ctx, qConfig.Name)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(q.Subscribers))
	for subName := range q.Subscribers {
		assert.Equal(t, subscriberName, subName)
	}
}

func TestQueuesUC_CantSubscribeToFullQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName1 := "subscriber1"
	subscriberName2 := "subscriber2"

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName1)
	assert.Nil(t, err)
	err = queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName2)
	assert.NotNil(t, err)
}

func TestQueuesUC_SubscriberCanConsumeMessagesFromQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName := "subscriber"
	msgBody := map[string]interface{}{"msg": "hello"}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName)
	assert.Nil(t, err)

	err = queuesUC.AddMessage(ctx, qConfig.Name, msgBody)
	assert.Nil(t, err)

	messsages, err := queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(messsages))
	for _, mes := range messsages {
		assert.Equal(t, msgBody, mes)
	}

	messsages, err = queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messsages))
}

func TestQueuesUC_NonSubscriberCantConsumeMessagesFromQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 1,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName := "not_a_subscriber"
	msgBody := map[string]interface{}{"msg": "hello"}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddMessage(ctx, qConfig.Name, msgBody)
	assert.Nil(t, err)

	_, err = queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName)
	assert.NotNil(t, err)
}

func TestQueuesUC_SubscribersCanConsumeMessagesFromQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 2,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName1 := "subscriber 1"
	subscriberName2 := "subscriber 2"
	msgBody := map[string]interface{}{"msg": "hello"}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName1)
	assert.Nil(t, err)
	err = queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName2)
	assert.Nil(t, err)

	err = queuesUC.AddMessage(ctx, qConfig.Name, msgBody)
	assert.Nil(t, err)

	messsages1, err := queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messsages1))
	for _, mes := range messsages1 {
		assert.Equal(t, msgBody, mes)
	}

	messsages2, err := queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName2)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messsages2))
	for _, mes := range messsages2 {
		assert.Equal(t, msgBody, mes)
	}

	messsages1_1, err := queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messsages1_1))
	messsages2_1, err := queuesUC.ConsumeMessages(ctx, qConfig.Name, subscriberName1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messsages2_1))
}

func TestQueuesUC_CantAddSubscribersToNotExistentQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName := "subscriber"

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, "some name", subscriberName)
	assert.NotNil(t, err)
}

func TestQueuesUC_CantAddMessagesToNotExistentQueue(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	msgBody := map[string]interface{}{"msg": "hello"}

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddMessage(ctx, "some name", msgBody)
	assert.NotNil(t, err)
}

func TestQueuesUC_CantAddSubscribersToQueueTwice(t *testing.T) {
	t.Parallel()

	qConfig := config.QueueConfig{
		Name:              "testQueue",
		Length:            1,
		SubscribersAmount: 2,
	}

	ctrl := gomock.NewController(t)
	qs := []config.QueueConfig{
		qConfig,
	}

	subscriberName := "subscriber"

	queuesUC, cleanup := configureEnvironment(ctrl, qs)
	defer cleanup()

	ctx := context.Background()

	err := queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName)
	assert.Nil(t, err)

	err = queuesUC.AddSubscriber(ctx, qConfig.Name, subscriberName)
	assert.NotNil(t, err)
}
