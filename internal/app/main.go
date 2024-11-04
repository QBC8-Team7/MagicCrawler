package app

import (
	"log"
	"time"

	"gopkg.in/telebot.v4"
)

type BotServer struct {
	Bot *telebot.Bot
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

	bot.Handle("/hello", func(c telebot.Context) error {
		return c.Send("Helloooo!")
	})

	bot.Start()
}
