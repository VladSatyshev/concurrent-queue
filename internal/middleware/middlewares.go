package middleware

import (
	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/queues"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
)

type MiddlewareManager struct {
	queuesUC queues.UseCase
	cfg      *config.Config
	logger   logger.Logger
}

func NewMiddlewareManager(queuesUC queues.UseCase, cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		queuesUC: queuesUC,
		cfg:      cfg,
		logger:   logger,
	}
}
