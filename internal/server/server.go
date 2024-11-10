package server

import (
	"context"
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
}

type UserContext struct {
	Command   string
	CurrentAd *Ad
	Progress  int
}

var userContext = make(map[int64]*UserContext)

func NewServer(ctx context.Context, cfg *config.Config, db *sqlc.Queries) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	handler := &Handlers{
		Logger: appLogger,
		DB:     db,
		DbCtx:  ctx,
	}
	log.Printf("Authorized on account %s\n\n", bot.Self.UserName)

	bot.Debug = cfg.Server.Mode == "development"

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
			case "/start":
				sendWellcome(bot, chatID, update.Message.From)
			default:
				handleUserMessage(ctx, bot, update, chatID, *db)
			}
		}
		if update.CallbackQuery != nil {
			action := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID

			switch action {
			case "ad_create":
				userContext[chatID] = &UserContext{
					Command:   "addhouse",
					CurrentAd: &Ad{},
					Progress:  0,
				}
				sendCategoryButtons(bot, chatID)
			case "ad_update":
				userContext[chatID] = &UserContext{
					Command:   "updatehouse",
					CurrentAd: &Ad{},
					Progress:  0,
				}
				sendCategoryButtons(bot, chatID)
			}

			handleCallbackQuery(bot, update)
		}
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
	}
}
