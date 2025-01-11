package middleware

import (
	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
)

type MiddlewareManager struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:    cfg,
		logger: logger,
	}
}
