package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/internal/middleware"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Server struct {
	router    *echo.Echo
	logger    *logger.AppLogger
	cfg       *config.Config
	db        *sqlc.Queries
	dbContext context.Context
}

func NewServer(dbCtx context.Context, cfg *config.Config, db *sqlc.Queries) (*Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	e := echo.New()

	s := &Server{
		router:    e,
		logger:    appLogger,
		cfg:       cfg,
		db:        db,
		dbContext: dbCtx,
	}

	registerRoutes(e, s)
	return s, nil
}

func (s *Server) Run() error {
	// TODO do we need to keep the simple bot server?
	//go func() {
	//	bot, err := tgbot.NewBotAPI(s.cfg.Bot.Token)
	//	if err != nil {
	//		log.Fatalf("Telegram Bot API initialization error: %v", err)
	//	}
	//	log.Println("Telegram Bot API initialized", bot.Self.ID)
	//
	//	http.HandleFunc("/bot", CreateBotEndpointHandler(bot, "https://6926-178-63-176-230.ngrok-free.app/"))
	//
	//	log.Fatal(http.ListenAndServe(s.cfg.Server.Port, nil))
	//}()

	s.router.Use(echoMiddlewares.CORSWithConfig(echoMiddlewares.CORSConfig{
		AllowOrigins: []string{s.cfg.Server.Host, "https://qbc8.boloorin.top.com", "https://6926-178-63-176-230.ngrok-free.app/"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	s.router.Use(middleware.EchoRequestLogger(s.logger))
	s.router.Use(middleware.EchoAuthentication(s.dbContext, s.db))
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	log.Println(addr)
	return s.router.Start(addr)
}

const SourceCodeUrl = "https://github.com/your-repo" // Define the source code URL

// CreateBotEndpointHandler : According to the https://core.telegram.org/bots/api#setwebhook webhook will receive JSON-serialized Update structure
// Handler created by this function parses Update structure and replies to any message with welcome text and inline keyboard to open Mini App
func CreateBotEndpointHandler(bot *tgbot.BotAPI, appURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving %s route", r.URL.Path)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		var update tgbot.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if update.Message == nil {
			http.Error(w, "Bot update didn't include a message", http.StatusBadRequest)
			return
		}

		message := "Welcome to the Telegram Mini App Template Bot\nTap the button below to open mini app or bot source code"
		inlineKeyboard := tgbot.NewInlineKeyboardMarkup(
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonData("Open mini app", appURL),
			),
			tgbot.NewInlineKeyboardRow(
				tgbot.NewInlineKeyboardButtonURL("Open source code", SourceCodeUrl),
			),
		)

		msg := tgbot.NewMessage(update.Message.Chat.ID, message)
		msg.ReplyMarkup = inlineKeyboard

		if _, err := bot.Send(msg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
