package server

import "gopkg.in/telebot.v4"

func GenerateRoutes(bot *telebot.Bot, h Handlers) {
	bot.Handle("/hello", h.HandleHello)
	bot.Handle("/bye", h.HandleBye)
}
