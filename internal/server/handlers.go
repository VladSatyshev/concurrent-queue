package server

import (
	"context"

	"github.com/VladSatyshev/concurrent-queue/internal/middleware"
	queuesHttp "github.com/VladSatyshev/concurrent-queue/internal/queues/delivery/http"
	queuesRepo "github.com/VladSatyshev/concurrent-queue/internal/queues/repository"
	queuesUseCase "github.com/VladSatyshev/concurrent-queue/internal/queues/usecase"
)

func (s *Server) MapHandlers() error {
	// init repositories
	qRepo := queuesRepo.NewQueuesRepository(s.cfg)
	if err := queuesRepo.InitQueues(context.Background(), s.cfg, qRepo); err != nil {
		s.logger.Errorf("failed to init queues: %s", err.Error())
		return err
	}

	// init usecases
	queuesUC := queuesUseCase.NewQueuesUseCase(s.cfg, qRepo, s.logger)

	// init handlers
	queuesHandlers := queuesHttp.NewQueuesHndlers(s.cfg, queuesUC, s.logger)

	// init & use middleware
	mw := middleware.NewMiddlewareManager(s.cfg, s.logger)
	s.router.Use(mw.CORSMiddleware())
	s.router.Use(mw.TimeoutMiddleware())

	v1 := s.router.Group("/v1")
	internal := v1.Group("/int")

	queueGroup := v1.Group("/queues")
	intQueueGroup := internal.Group("/queues")

	queuesHttp.MapQueueRoutes(queueGroup, queuesHandlers, mw)
	queuesHttp.MapIntQueueRoutes(intQueueGroup, queuesHandlers, mw)

	return nil
}
