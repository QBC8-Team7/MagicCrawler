package server

import (
	"log"
	"strconv"

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

			switch text {
			case "/addhouse":
				userContext[userID] = &UserContext{
					Command:   "addhouse",
					CurrentAd: &Ad{},
					Progress:  0,
				}
				sendCategoryButtons(bot, userID)

			case "/updatehouse":
				userContext[userID] = &UserContext{
					Command:   "updatehouse",
					CurrentAd: &Ad{}, // TODO: Load the ad to be updated here
					Progress:  0,
				}
				sendCategoryButtons(bot, userID)
			default:
				handleUserMessage(bot, update, userID)
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

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.CallbackQuery.Message.Chat.ID
	context, ok := userContext[userID]
	if !ok {
		return
	}

	ad := context.CurrentAd
	switch context.Progress {
	case 0: // Category
		ad.Category = update.CallbackQuery.Data
		sendHouseTypeButtons(bot, userID)
		context.Progress++

	case 1: // HouseType
		ad.Category = update.CallbackQuery.Data
		sendWarehouseButtons(bot, userID)
		context.Progress++

	case 2: // Warehouse
		ad.Category = update.CallbackQuery.Data
		sendElevatorButtons(bot, userID)
		context.Progress++

	case 3: // Elevator
		ad.Category = update.CallbackQuery.Data
		bot.Send(tgbotapi.NewMessage(userID, "Enter Publisher Ad Key"))
		context.Progress++
	}

	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

func sendCategoryButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Rent", "rent"),
			tgbotapi.NewInlineKeyboardButtonData("Buy", "buy"),
			tgbotapi.NewInlineKeyboardButtonData("Mortgage", "mortgage"),
		),
	)
	msg := tgbotapi.NewMessage(userID, "Select Ad Category")
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func sendHouseTypeButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Apartment", "apartment"),
			tgbotapi.NewInlineKeyboardButtonData("Villa", "villa"),
		),
	)
	msg := tgbotapi.NewMessage(userID, "Select House Type")
	msg.ReplyMarkup = buttons
	bot.Send(msg)
}

func sendWarehouseButtons(bot *tgbotapi.BotAPI, userID int64) {
	// Define inline buttons for HasWarehouse
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		),
	)

	msg := tgbotapi.NewMessage(userID, "Does the house have a warehouse?")
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

func handleUserMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, userID int64) {
	context, inProgress := userContext[userID]
	if !inProgress {
		bot.Send(tgbotapi.NewMessage(userID, "Use /addhouse or /updatehouse to start."))
		return
	}

	ad := context.CurrentAd
	text := update.Message.Text
	switch context.Progress {
	case 4:
		// TODO: validate publisher ad key
		if text != "" {
			ad.PublisherAdKey = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Publisher ID"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Publisher Ad Key again"))
		}
	case 5: // PublisherID
		// TODO: validation
		pid, err := strconv.Atoi(text)
		if err == nil && pid > 0 {
			ad.PublisherID = pid
			bot.Send(tgbotapi.NewMessage(userID, "Enter Author"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Publisher ID again"))
		}
	case 6: // Author
		// TODO: validation
		if text != "" {
			ad.Author = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Title"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Author again"))
		}
	case 7: // Title
		// TODO: validation
		if text != "" {
			ad.Title = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Description"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Title again"))
		}
	case 8: // Description
		// TODO: validation
		if text != "" {
			ad.Description = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter City"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Description again"))
		}
	case 9: // City
		// TODO: validation
		if text != "" {
			ad.City = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Neighborhood"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter City again"))
		}
	case 10: // Neighborhood
		// TODO: validation
		if text != "" {
			ad.Neighborhood = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Meterage"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Neighborhood again"))
		}
	case 11: // Meterage
		// TODO: validation
		meterage, err := strconv.Atoi(text)
		if err == nil && meterage >= 0 {
			ad.Meterage = meterage
			bot.Send(tgbotapi.NewMessage(userID, "Enter Rooms Count"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Meterage again"))
		}
	case 12: // RoomsCount
		// TODO: validation
		rooms, err := strconv.Atoi(text)
		if err == nil && rooms >= 0 {
			ad.RoomsCount = rooms
			bot.Send(tgbotapi.NewMessage(userID, "Enter Year"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Rooms Count again"))
		}

	case 13: // Year
		// TODO: validation
		year, err := strconv.Atoi(text)
		if err == nil && year >= 1250 {
			ad.Year = year
			bot.Send(tgbotapi.NewMessage(userID, "Enter Floor number"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter year again"))
		}

	case 14: // Floor
		// TODO: validation
		floor, err := strconv.Atoi(text)
		if err == nil && floor >= 0 {
			ad.Floor = floor
			bot.Send(tgbotapi.NewMessage(userID, "Enter Total Floors"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter floor again"))
		}

	case 15: // Total Floors
		// TODO: validation
		totalFloors, err := strconv.Atoi(text)
		if err == nil && totalFloors >= 0 {
			ad.TotalFloors = totalFloors
			bot.Send(tgbotapi.NewMessage(userID, "Enter house latitude"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter total floors again"))
		}

	case 16: // Lat
		// TODO: validation
		if text != "" {
			ad.Lat = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter house longitude"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter latitude again"))
		}

	case 17: // Lon
		// TODO: validation
		if text != "" {
			ad.Lng = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter ad URL"))
			context.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter longitude again"))
		}

	case 18: // URL
		// TODO: validation
		if text != "" {
			ad.Url = text
			context.Progress++

			if isAdComplete(ad) {
				if context.Command == "addhouse" {
					// TODO: insert ad to DB
					bot.Send(tgbotapi.NewMessage(userID, "House added successfully."))
				} else if context.Command == "updatehouse" {
					// TODO: update ad
					bot.Send(tgbotapi.NewMessage(userID, "House updated successfully."))
				}
				delete(userContext, userID)
			}

		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter URL again"))
		}
	}

}

func isAdComplete(ad *Ad) bool {
	// TODO: check if ad is OK
	return true
}
