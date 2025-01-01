package server

import (
	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
	logger logger.Logger
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,

		router: gin.New(),
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(); err != nil {
		return err
	}

	if err := s.router.Run(s.cfg.Server.Port); err != nil {
		return err
	}

	return nil
}
