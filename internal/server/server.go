package server

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"log"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/middleware"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

type BotServer struct {
	Bot     *telebot.Bot
	Handler *Handler
	Logger  *logger.AppLogger
	DB      *sqlc.Queries
}

func NewServer(cfg *config.Config, db *sqlc.Queries) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	botSetting := telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(botSetting)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	handler := &Handler{
		Logger: appLogger,
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
		DB:      db,
	}
}

func (s *BotServer) Serve() {
	s.Bot.Use(middleware.WithLogging(s.Logger))

	GenerateRoutes(s)
	s.Bot.Start()
}
