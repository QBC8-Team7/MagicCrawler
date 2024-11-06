package server

import (
	"log"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

type BotServer struct {
	Bot     *telebot.Bot
	Handler *Handlers
	Logger  *logger.AppLogger
}

func NewServer(cfg *config.Config) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger("app.log")
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	settings := telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	handler := &Handlers{
		Logger: appLogger,
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
	}
}

func (s *BotServer) Serve() {
	GenerateRoutes(s)

	s.Bot.Start()
}
