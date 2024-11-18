package notification

import (
	"strconv"
	"sync"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Service struct {
	Bot    *tgbotapi.BotAPI
	Logger *logger.AppLogger
}

var (
	instance *Service
	once     sync.Once
)

func GetService(botToken string, appLogger *logger.AppLogger) (*Service, error) {
	var initErr error
	once.Do(func() {
		bot, err := tgbotapi.NewBotAPI(botToken)
		if err != nil {
			initErr = err
			return
		}
		appLogger.Infof("Telegram bot authorized on account %s", bot.Self.UserName)
		instance = &Service{Bot: bot, Logger: appLogger}
	})
	return instance, initErr
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
