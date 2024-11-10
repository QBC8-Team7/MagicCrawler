package server

import (
	"log"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"

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

type UserContext struct {
	Command   string
	CurrentAd *Ad
	Progress  int
}

var userContext = make(map[int64]*UserContext)

func NewServer(cfg *config.Config) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	handler := &Handler{
		Logger: appLogger,
	}
	// bot.Debug = cfg.Server.Mode == "development"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			switch text {
			case "/addhouse":
				userContext[chatID] = &UserContext{
					Command:   "addhouse",
					CurrentAd: &Ad{},
					Progress:  0,
				}
				sendCategoryButtons(bot, chatID)

			case "/updatehouse":
				userContext[chatID] = &UserContext{
					Command:   "updatehouse",
					CurrentAd: &Ad{}, // TODO: Load the ad to be updated here
					Progress:  0,
				}
				sendCategoryButtons(bot, chatID)
			default:
				handleUserMessage(bot, update, chatID)
			}
		}

		if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update)
		}
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
		DB:      db,
	}
}
