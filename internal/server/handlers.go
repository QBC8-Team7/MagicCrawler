package server

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

type CommandHandler interface {
	HandleHello(c telebot.Context) error
	HandleBye(c telebot.Context) error
}

type Handler struct {
	Logger *logger.AppLogger
}

func (h *Handler) HandleHello(c telebot.Context) error {
	h.Logger.Info("log from hello")
	return c.Send("Helloooo!")
}

func (h *Handler) HandleBye(c telebot.Context) error {
	h.Logger.Info("log from Bye")

	return c.Send("Bye!")
}
