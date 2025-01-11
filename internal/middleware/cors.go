package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) CORSMiddleware() gin.HandlerFunc {
	mw.logger.Info("Setting CORS")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "X-Subscriber")
	return cors.New(config)
}
