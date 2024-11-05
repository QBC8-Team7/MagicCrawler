package main

import (
	"log"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/app"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger("app.log")
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	s := app.NewServer(cfg.Token)
	s.Serve()
}
