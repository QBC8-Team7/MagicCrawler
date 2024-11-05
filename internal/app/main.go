package app

import (
	"log"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/internal/server"
	handlers "github.com/QBC8-Team7/MagicCrawler/internal/server"
	"gopkg.in/telebot.v4"
)

type BotServer struct {
	Bot     *telebot.Bot
	Handler handlers.Handlers
}

func NewServer(token string) *BotServer {
	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	return &BotServer{Bot: bot}
}

func (s *BotServer) Serve() {
	bot := s.Bot

	server.GenerateRoutes(bot, s.Handler)

	bot.Start()
}
