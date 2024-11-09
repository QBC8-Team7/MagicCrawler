package server

import (
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"log"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotServer struct {
	Bot     *tgbotapi.BotAPI
	Handler *Handlers
	Logger  *logger.AppLogger
	DB      *sqlc.Queries
}

func NewServer(cfg *config.Config, db *sqlc.Queries) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

<<<<<<< HEAD
	botSetting := telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(botSetting)
=======
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
>>>>>>> bad6448 (feat: intial setup with tgbot)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	handler := &Handler{
		Logger: appLogger,
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil { // If there is a message in this update
			// Call the appropriate handler based on the message content
			switch update.Message.Text {
			case "/start":
				// Use your custom handler for the /start command
				err := handler.HandleStart(update.Message, bot)
				if err != nil {
					log.Printf("Error handling /start: %v", err)
				}
			case "/help":
				// Send a help message as an example (you can create a handler for this if needed)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Show help")
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			default:
				// Default behavior for unknown commands
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /help for available commands.")
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Failed to send message: %v", err)
				}
			}
		}
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
		DB:      db,
	}
}
