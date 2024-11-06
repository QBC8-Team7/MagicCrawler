package middleware

import (
	"log"
	"time"

	"gopkg.in/telebot.v4"
)

func LoggingMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {

	return func(c telebot.Context) error {
		msg := c.Message()
		if msg != nil {
			log.Printf("Received message from %s (%d): %s at %s",
				msg.Sender.Username,
				msg.Sender.ID,
				msg.Text,
				time.Now().Format(time.RFC3339))
		}
		return next(c)
	}
}
