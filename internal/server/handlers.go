package server

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

type CommandHandler interface {
	HandleHello(c telebot.Context) error
	HandleBye(c telebot.Context) error
}

type Handlers struct {
	Logger *logger.AppLogger
}

func (h *Handlers) HandleHello(c telebot.Context) error {
	h.Logger.Info("log from hello")
	return c.Send("Helloooo!")
}

func (h *Handlers) HandleBye(c telebot.Context) error {
	h.Logger.Info("log from Bye")

	return c.Send("Bye!")
}
