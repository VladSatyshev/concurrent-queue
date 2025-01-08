package middleware

import (
	"net/http"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func (mw *MiddlewareManager) TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(mw.cfg.Server.Timeout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.String(http.StatusRequestTimeout, "timeout")
		}),
	)
}
