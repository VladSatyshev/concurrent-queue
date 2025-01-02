package http

import (
	"fmt"
	"net/http"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
	"github.com/VladSatyshev/concurrent-queue/pkg/utils"
	"github.com/gin-gonic/gin"
)

type queuesHandlers struct {
	cfg      *config.Config
	queuesUC queues.UseCase
	logger   logger.Logger
}

func NewQueuesHndlers(cfg *config.Config, queuesUC queues.UseCase, log logger.Logger) queues.Handlers {
	return &queuesHandlers{cfg: cfg, queuesUC: queuesUC, logger: log}
}

func handleError(c *gin.Context, err error) {
	if qErr, ok := err.(*queues.QueueErr); !ok {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		switch qErr.ErrType {
		case queues.RepositoryErr, queues.RepositoryNotFoundErr:
			c.JSON(http.StatusInternalServerError, err.Error())
		case queues.UseCaseNotFoundErr:
			c.JSON(http.StatusNotFound, err.Error())
		case queues.UseCaseErr:
			c.JSON(http.StatusBadRequest, err.Error())
		}
	}
}

func (h *queuesHandlers) GetAll() func(c *gin.Context) {
	return func(c *gin.Context) {
		queues := h.queuesUC.GetAll(c.Request.Context())
		c.JSON(http.StatusOK, queues)
	}
}

func (h *queuesHandlers) GetQueueByName() func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("queue_name")

		queue, err := h.queuesUC.GetByName(c.Request.Context(), name)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, queue)
	}
}

func (h *queuesHandlers) Subscribe() func(c *gin.Context) {
	return func(c *gin.Context) {

		queueName := c.Param("queue_name")

		subscriberName, err := utils.GetSubscriber(c)
		if err != nil {
			handleError(c, err)
			return
		}

		err = h.queuesUC.AddSubscriber(c.Request.Context(), queueName, subscriberName)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("user %v has subscribed to queue %s", subscriberName, queueName))
	}
}

func (h *queuesHandlers) AddMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		queueName := c.Param("queue_name")

		var jsonBody map[string]interface{}
		if err := c.Bind(&jsonBody); err != nil {
			h.logger.Errorf("failed to parse json body: %s", err.Error())
			handleError(c, err)
			return
		}

		err := h.queuesUC.AddMessage(c.Request.Context(), queueName, jsonBody)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("message has been added to queue %s", queueName))
	}
}

func (h *queuesHandlers) Consume() func(c *gin.Context) {
	return func(c *gin.Context) {
		queueName := c.Param("queue_name")

		subscriberName, err := utils.GetSubscriber(c)
		if err != nil {
			handleError(c, err)
			return
		}

		messages, err := h.queuesUC.ConsumeMessages(c.Request.Context(), queueName, subscriberName)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}
