package main

import (
	"flag"
	"log"

	"github.com/VladSatyshev/concurrent-queue/config"
	"github.com/VladSatyshev/concurrent-queue/internal/server"
	"github.com/VladSatyshev/concurrent-queue/pkg/logger"
	"github.com/VladSatyshev/concurrent-queue/pkg/utils"
)

const defaultConfigPath = "./config/config-local.yml"

func main() {
	log.Println("Starting API server")

	configPathFlag := flag.String("config", defaultConfigPath, "config path")
	flag.Parse()

	configPath := utils.GetConfigPath(*configPathFlag)

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Server.Mode)

	s := server.NewServer(cfg, appLogger)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
