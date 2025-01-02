package queues

import "github.com/gin-gonic/gin"

type Handlers interface {
	// int
	GetAll() func(*gin.Context)
	GetQueueByName() func(*gin.Context)

	// public
	Subscribe() func(*gin.Context)
	AddMessage() func(*gin.Context)
	Consume() func(*gin.Context)
}
