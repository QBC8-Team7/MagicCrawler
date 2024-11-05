package server

import "gopkg.in/telebot.v4"

type CommandHandler interface {
	HandleHello(c telebot.Context) error
	HandleBye(c telebot.Context) error
}

type Handlers struct{}

func (h *Handlers) HandleHello(c telebot.Context) error {
	return c.Send("Helloooo!")
}

func (h *Handlers) HandleBye(c telebot.Context) error {
	return c.Send("Bye!")
}
