package notification

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type Service struct {
	Bot    *tgbotapi.BotAPI
	Logger *logger.AppLogger
}

func NewNotificationService(botToken string, logger *logger.AppLogger) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	logger.Infof("Telegram bot uthorized on account %s", bot.Self.UserName)

	return &Service{Bot: bot, Logger: logger}, nil
}

func (n *Service) SendMessage(tgID string, message string) error {
	userID, err := strconv.ParseInt(tgID, 10, 64)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, message)
	_, err = n.Bot.Send(msg)
	if err != nil {
		n.Logger.Errorf("Failed to send message to user %d: %v", userID, err)
		return err
	}
	n.Logger.Infof("Message sent to user %d: %s", userID, message)
	return nil
}
