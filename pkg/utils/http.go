package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetConfigPath(configPath string) string {
	switch configPath {
	default:
		return "./config/config-local.yml"
	}
}

func GetSubscriber(c *gin.Context) (string, error) {
	subscriberName, ok := c.Request.Header["X-Subscriber"]
	if !ok {
		return "", errors.New("failed to parse X-Subscriber header")
	}
	if len(subscriberName) != 1 {
		return "", errors.New("only one subscriber in X-Subscriber header is allowed")
	}

	return subscriberName[0], nil
}
