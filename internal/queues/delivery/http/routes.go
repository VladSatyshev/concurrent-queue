package http

import (
	"github.com/VladSatyshev/concurrent-queue/internal/middleware"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/gin-gonic/gin"
)

func MapIntQueueRoutes(intQueueGroup *gin.RouterGroup, h queues.Handlers, mw *middleware.MiddlewareManager) {
	intQueueGroup.GET("/", h.GetAll())
	intQueueGroup.GET("/:queue_name", h.GetQueueByName())
}

func MapQueueRoutes(queueGroup *gin.RouterGroup, h queues.Handlers, mw *middleware.MiddlewareManager) {
	queueGroup.POST("/:queue_name/subscriptions", h.Subscribe())
	queueGroup.POST("/:queue_name/messages", h.AddMessage())
}
