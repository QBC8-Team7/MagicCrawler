package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotServer struct {
	Bot     *tgbotapi.BotAPI
	Handler *Handlers
	Logger  *logger.AppLogger
}

var userStates = make(map[int64]*Ad)
var userProgress = make(map[int64]int)

func NewServer(cfg *config.Config) *BotServer {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	handler := &Handlers{
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
		if update.Message != nil {
			userID := update.Message.Chat.ID
			text := update.Message.Text

			if text == "/addhouse" {
				// Initialize ad and progress
				userStates[userID] = &Ad{}
				userProgress[userID] = 0
				bot.Send(tgbotapi.NewMessage(userID, "Enter Publisher Ad Key:"))
				continue
			}

			// Get user state and progress
			ad, inProgress := userStates[userID]
			progress, hasProgress := userProgress[userID]
			if !inProgress || !hasProgress {
				bot.Send(tgbotapi.NewMessage(userID, "Use /addhouse to start adding a new house."))
				continue
			}

			// Handle each field by progress index
			switch progress {
			case 0: // PublisherAdKey
				ad.PublisherAdKey = text
				bot.Send(tgbotapi.NewMessage(userID, ""))
				userProgress[userID]++
			case 1: // PublisherID
				if text != "" {
					pid, err := strconv.Atoi(text)
					if err == nil {
						ad.PublisherID = pid
					} else {
						bot.Send(tgbotapi.NewMessage(userID, "Invalid number. Enter Publisher ID:"))
						continue
					}
				}
				bot.Send(tgbotapi.NewMessage(userID, "Enter Category (buy, rent, mortgage):"))
				userProgress[userID]++
			case 2: // Category
				ad.Category = text
				bot.Send(tgbotapi.NewMessage(userID, "Does the house have an elevator?"))
				sendElevatorButtons(bot, userID)
				userProgress[userID]++
			case 12: // Author (optional)
				if text != "" {
					ad.Author = text
				}
				bot.Send(tgbotapi.NewMessage(userID, "Enter URL (optional, press Enter to skip):"))
				userProgress[userID]++
			// Continue handling other fields similarly...
			case 4: // Title
				ad.Title = text
				bot.Send(tgbotapi.NewMessage(userID, "Enter Description (optional, press Enter to skip):"))
				userProgress[userID]++
			case 5: // Description
				if text != "" {
					ad.Description = text
				}
				bot.Send(tgbotapi.NewMessage(userID, "Enter City:"))
				userProgress[userID]++
			case 6: // City
				ad.City = text
				bot.Send(tgbotapi.NewMessage(userID, "Enter Neighborhood (optional, press Enter to skip):"))
				userProgress[userID]++
			case 7: // Neighborhood
				if text != "" {
					ad.Neighborhood = text
				}
				bot.Send(tgbotapi.NewMessage(userID, "Enter Meterage:"))
				userProgress[userID]++
			case 8: // Meterage
				meterage, err := strconv.Atoi(text)
				if err == nil && meterage >= 0 {
					ad.Meterage = meterage
					bot.Send(tgbotapi.NewMessage(userID, "Enter Rooms Count:"))
					userProgress[userID]++
				} else {
					bot.Send(tgbotapi.NewMessage(userID, "Invalid number. Enter Meterage:"))
				}
			case 9: // RoomsCount
				rooms, err := strconv.Atoi(text)
				if err == nil && rooms >= 0 {
					ad.RoomsCount = rooms
					bot.Send(tgbotapi.NewMessage(userID, "Enter Year:"))
					userProgress[userID]++
				} else {
					bot.Send(tgbotapi.NewMessage(userID, "Invalid number. Enter Rooms Count:"))
				}

			case 10: // Year
				year, err := strconv.Atoi(text)
				if err == nil && year >= 1250 {
					ad.Year = year
					bot.Send(tgbotapi.NewMessage(userID, "Enter Floor number:"))
					userProgress[userID]++
				} else {
					bot.Send(tgbotapi.NewMessage(userID, "Invalid year. Enter year"))
				}

			case 11: // Floor
				floor, err := strconv.Atoi(text)
				if err == nil && floor >= 0 {
					ad.Floor = floor
					bot.Send(tgbotapi.NewMessage(userID, "Enter Total Floors:"))
					userProgress[userID]++
				} else {
					bot.Send(tgbotapi.NewMessage(userID, "Invalid floor. Enter floor"))
				}

				// case 12: // Floor
				// 	totalFloors, err := strconv.Atoi(text)
				// 	if err == nil && totalFloors >= 0 {
				// 		ad.TotalFloors = totalFloors
				// 		bot.Send(tgbotapi.NewMessage(userID, "HasWarehouse: (y/n)"))
				// 		userProgress[userID]++
				// 	} else {
				// 		bot.Send(tgbotapi.NewMessage(userID, "Invalid floor. Enter floor"))
				// 	}

				// Floor          int
				// TotalFloors    int
				// HasWarehouse   bool
				// HasElevator    bool
				// Lat            float64
				// Lng            float64

				// Continue with additional fields like Floor, TotalFloors, HasWarehouse, HasElevator, Lat, Lng, etc.
			}

			// After all fields are filled
			if isAdComplete(ad) {
				insertHouseToDB(ad)
				bot.Send(tgbotapi.NewMessage(userID, "House added successfully!"))
				// Clean up user state
				delete(userStates, userID)
				delete(userProgress, userID)
			}
		}

		if update.CallbackQuery != nil {
			userID := update.CallbackQuery.Message.Chat.ID
			data := update.CallbackQuery.Data

			if ad, ok := userStates[userID]; ok {
				progress := userProgress[userID]

				switch progress {
				case 3: // HasElevator
					ad.HasElevator = data == "yes"
					bot.Send(tgbotapi.NewMessage(userID, "Thank you! House data saved."))
					// Insert ad into DB here if all fields are filled
					delete(userStates, userID)
					delete(userProgress, userID)
				}

				// Answer callback to remove loading indicator
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			}
		}
	}

	return &BotServer{
		Bot:     bot,
		Handler: handler,
		Logger:  appLogger,
	}
}

func sendCategoryButtons(bot *tgbotapi.BotAPI, userID int64) {
	// Define inline buttons for Category
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Apartment", "apartment"),
			tgbotapi.NewInlineKeyboardButtonData("House", "house"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Villa", "villa"),
			tgbotapi.NewInlineKeyboardButtonData("Studio", "studio"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Select Category:")
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func sendElevatorButtons(bot *tgbotapi.BotAPI, userID int64) {
	// Define inline buttons for HasElevator
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Does the house have an elevator?")
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func allFieldsFilled(ad *Ad) bool {
	// Check if all required fields are filled
	return ad.PublisherAdKey != "" && ad.Category != "" // continue for other fields
}

func insertHouseToDB(ad *Ad) {
	// Insert ad into DB
	fmt.Printf("Inserting house into DB: %+v\n", ad)
	// Your database code here
}

func isAdComplete(ad *Ad) bool {
	return ad.PublisherAdKey != "" && ad.Category != "" && ad.Meterage > 0 && ad.RoomsCount > 0 // and other mandatory fields
}
