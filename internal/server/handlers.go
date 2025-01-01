package server

import (
	"github.com/VladSatyshev/concurrent-queue/internal/middleware"
	queuesHttp "github.com/VladSatyshev/concurrent-queue/internal/queues/delivery/http"
	queuesRepo "github.com/VladSatyshev/concurrent-queue/internal/queues/repository"
	queuesUseCase "github.com/VladSatyshev/concurrent-queue/internal/queues/usecase"
)

func (s *Server) MapHandlers() error {
	// init repositories
	queuesRepo, err := queuesRepo.NewQueuesRepository(s.cfg)
	if err != nil {
		s.logger.Errorf("failed to create queues repository: %s", err.Error())
		return err
	}

	// init usecases
	queuesUC := queuesUseCase.NewQueuesUseCase(s.cfg, queuesRepo, s.logger)

	// init handlers
	queuesHandlers := queuesHttp.NewQueuesHndlers(s.cfg, queuesUC, s.logger)

	// init & use middleware
	mw := middleware.NewMiddlewareManager(queuesUC, s.cfg, s.logger)
	s.router.Use(mw.TimeoutMiddleware())

	v1 := s.router.Group("/v1")
	internal := v1.Group("/int")

	queueGroup := v1.Group("/queues")
	intQueueGroup := internal.Group("/queues")

	queuesHttp.MapQueueRoutes(queueGroup, queuesHandlers, mw)
	queuesHttp.MapIntQueueRoutes(intQueueGroup, queuesHandlers, mw)

	return nil
}
