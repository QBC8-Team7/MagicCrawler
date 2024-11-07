package middleware

import (
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

func WithLogging(logger *logger.AppLogger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			msg := c.Message()
			if msg != nil {
				logger.Infof("Received message from %s (%d): %s at %s",
					msg.Sender.Username,
					msg.Sender.ID,
					msg.Text,
					time.Now().Format(time.RFC3339))
			}
			return next(c)
		}
	}
}
