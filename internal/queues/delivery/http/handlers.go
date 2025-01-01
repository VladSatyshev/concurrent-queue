package http

import (
	"net/http"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
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

func (h *queuesHandlers) GetAll() func(c *gin.Context) {
	return func(c *gin.Context) {
		queues, err := h.queuesUC.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, queues)
	}
}

func (h *queuesHandlers) GetQueueByName() func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("queue_name")

		queue, err := h.queuesUC.GetByName(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, queue)
	}
}

func (h *queuesHandlers) Subscribe() func(c *gin.Context) {
	return func(c *gin.Context) {

		queueName := c.Param("queue_name")

		subscriberName, ok := c.Request.Header["X-Subscriber"]
		if !ok {
			c.JSON(http.StatusBadRequest, "failed to parse X-Subscriber header")
		}
		if len(subscriberName) != 1 {
			c.JSON(http.StatusBadRequest, "only one subscriber in X-Subscriber header is allowed")
		}

		err := h.queuesUC.AddSubscriber(c.Request.Context(), queueName, subscriberName[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *queuesHandlers) AddMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		queueName := c.Param("queue_name")

		var jsonBody map[string]interface{}
		if err := c.BindJSON(jsonBody); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		h.queuesUC.AddMessage(c.Request.Context(), queueName, jsonBody)

		c.JSON(http.StatusOK, nil)
	}
}
