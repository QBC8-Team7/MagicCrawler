package server

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v5/pgtype"
)

type CommandHandler interface {
	HandleHello(update tgbotapi.Update) error
	HandleBye(update tgbotapi.Update) error
	HandleWatchlist(update tgbotapi.Update) error
}

type Handlers struct {
	Logger *logger.AppLogger
	DB     *sqlc.Queries
	DbCtx  context.Context
	Bot    *tgbotapi.BotAPI
	Ad     Ad
}

type Ad struct {
	PublisherAdKey string
	PublisherID    int
	Category       string
	Author         string
	Url            string
	Title          string
	Description    string
	City           string
	Neighborhood   string
	HouseType      string
	Meterage       int
	RoomsCount     int
	Year           int
	Floor          int
	TotalFloors    int
	HasWarehouse   bool
	HasElevator    bool
	Lat            string
	Lng            string
}

var lastMessageID = make(map[int64]int)

func replaceMessage(bot *tgbotapi.BotAPI, userID int64, text string, buttons *tgbotapi.InlineKeyboardMarkup) {
	if msgID, exists := lastMessageID[userID]; exists {
		editMsg := tgbotapi.NewEditMessageText(userID, msgID, text)
		editMsg.ReplyMarkup = buttons
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Println("Error editing message:", err)
			return
		}
	} else {
		msg := tgbotapi.NewMessage(userID, text)
		msg.ReplyMarkup = buttons

		sentMsg, err := bot.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		lastMessageID[userID] = sentMsg.MessageID
	}
}

func sendWellcome(bot *tgbotapi.BotAPI, userID int64, user *tgbotapi.User) {
	welcomeText := fmt.Sprintf("👋 Welcome, %s!\n\n", user.FirstName)
	welcomeText += "This bot helps you find home ads in Tehran 🏡. It gathers data from Shypoor and Divar 📱.\n\n"
	welcomeText += "What would you like to do today? 🤔\n\n"

	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Search Ads 🔍", "ad_search"),
			tgbotapi.NewInlineKeyboardButtonData("My Watch List 👀", "ad_watchlist"),
		),
	)

	// TODO: we must retirve user status(superadmin) from DB
	superUserID := int64(7417976949)
	if userID == superUserID {
		buttons.InlineKeyboard = append(buttons.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Create Ad 🛠️", "ad_create"),
				tgbotapi.NewInlineKeyboardButtonData("Update Ad 🔄", "ad_update"),
				tgbotapi.NewInlineKeyboardButtonData("Delete Ad 🗑️", "ad_delete"),
			),
		)
	}

	replaceMessage(bot, userID, welcomeText, &buttons)

}

func sendCategoryButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Rent", "rent"),
			tgbotapi.NewInlineKeyboardButtonData("Buy", "buy"),
			tgbotapi.NewInlineKeyboardButtonData("Mortgage", "mortgage"),
		),
	)
	replaceMessage(bot, userID, "Select Ad Category", &buttons)
}

func sendHouseTypeButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Apartment", "apartment"),
			tgbotapi.NewInlineKeyboardButtonData("Villa", "villa"),
		),
	)
	replaceMessage(bot, userID, "Select House Type", &buttons)
}

func sendWarehouseButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		),
	)
	replaceMessage(bot, userID, "Does the house have a warehouse?", &buttons)
}

func sendElevatorButtons(bot *tgbotapi.BotAPI, userID int64) {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("No", "no"),
		),
	)
	replaceMessage(bot, userID, "Does the house have an elevator?", &buttons)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.CallbackQuery.Message.Chat.ID
	userCtx, ok := userContext[userID]
	if !ok {
		return
	}

	ad := userCtx.CurrentAd
	switch userCtx.Progress {
	case 0: // Category
		// TODO: get categories from db
		ad.Category = update.CallbackQuery.Data
		sendHouseTypeButtons(bot, userID)
		userCtx.Progress++

	case 1: // HouseType
		// TODO: get house type from db
		ad.Category = update.CallbackQuery.Data
		sendWarehouseButtons(bot, userID)
		userCtx.Progress++

	case 2: // Warehouse
		ad.Category = update.CallbackQuery.Data
		sendElevatorButtons(bot, userID)
		userCtx.Progress++

	case 3: // Elevator
		ad.Category = update.CallbackQuery.Data
		bot.Send(tgbotapi.NewMessage(userID, "Enter Publisher Ad Key"))
		userCtx.Progress++
	}

	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

// TODO: pass logger
// TODO: support edit mode
func handleUserMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, userID int64, db sqlc.Queries) {
	userCtx, inProgress := userContext[userID]

	if !inProgress {
		bot.Send(tgbotapi.NewMessage(userID, "invalid command"))
		return
	}

	ad := userCtx.CurrentAd
	text := update.Message.Text

	switch userCtx.Progress {
	case 4:
		// TODO: validate publisher ad key
		// TODO: no need to ask in create mode. we must set Bot as default
		if text != "" {
			ad.PublisherAdKey = text

			bot.Send(tgbotapi.NewMessage(userID, "Enter Publisher ID"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Publisher Ad Key again"))
		}
	case 5: // PublisherID
		// TODO: validation
		pid, err := strconv.Atoi(text)
		if err == nil && pid > 0 {
			ad.PublisherID = pid
			bot.Send(tgbotapi.NewMessage(userID, "Enter Author"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Publisher ID again"))
		}
	case 6: // Author
		if text != "" {
			ad.Author = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Title"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Author again"))
		}
	case 7: // Title
		if text != "" {
			ad.Title = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Description"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Title again"))
		}
	case 8: // Description
		if text != "" {
			ad.Description = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter City"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Description again"))
		}
	case 9: // City
		if text != "" {
			ad.City = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Neighborhood"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter City again"))
		}
	case 10: // Neighborhood
		if text != "" {
			ad.Neighborhood = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter Meterage"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Neighborhood again"))
		}
	case 11: // Meterage
		meterage, err := strconv.Atoi(text)
		if err == nil && meterage >= 0 {
			ad.Meterage = meterage
			bot.Send(tgbotapi.NewMessage(userID, "Enter Rooms Count"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Meterage again"))
		}
	case 12: // RoomsCount
		rooms, err := strconv.Atoi(text)
		if err == nil && rooms >= 0 {
			ad.RoomsCount = rooms
			bot.Send(tgbotapi.NewMessage(userID, "Enter Year"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter Rooms Count again"))
		}

	case 13: // Year
		year, err := strconv.Atoi(text)
		if err == nil && year >= 1250 {
			ad.Year = year
			bot.Send(tgbotapi.NewMessage(userID, "Enter Floor number"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter year again"))
		}

	case 14: // Floor
		floor, err := strconv.Atoi(text)
		if err == nil && floor >= 0 {
			ad.Floor = floor
			bot.Send(tgbotapi.NewMessage(userID, "Enter Total Floors"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter floor again"))
		}

	case 15: // Total Floors
		totalFloors, err := strconv.Atoi(text)
		if err == nil && totalFloors >= 0 {
			ad.TotalFloors = totalFloors
			bot.Send(tgbotapi.NewMessage(userID, "Enter house latitude"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter total floors again"))
		}

	case 16: // Lat
		// TODO: validation
		if text != "" {
			ad.Lat = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter house longitude"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter latitude again"))
		}

	case 17: // Lon
		// TODO: validation
		if text != "" {
			ad.Lng = text
			bot.Send(tgbotapi.NewMessage(userID, "Enter ad URL"))
			userCtx.Progress++
		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter longitude again"))
		}

	case 18: // URL
		// TODO: url validation
		if text != "" {
			ad.Url = text
			userCtx.Progress++

			if userCtx.Command == "addhouse" {
				// TODO: insert ad to DB
				myAd := userCtx.CurrentAd

				ad := &sqlc.CreateAdParams{
					PublisherAdKey: myAd.PublisherAdKey,
					PublisherID:    pgtype.Int4{Int32: int32(myAd.PublisherID), Valid: true},
					PublishedAt:    pgtype.Timestamp{Time: time.Now(), Valid: true},
					Category:       sqlc.NullAdCategory{AdCategory: sqlc.AdCategory(myAd.Category), Valid: myAd.Category != ""},
					Author:         pgtype.Text{String: myAd.Author, Valid: myAd.Author != ""},
					Url:            pgtype.Text{String: myAd.Url, Valid: myAd.Url != ""},
					Title:          pgtype.Text{String: myAd.Title, Valid: myAd.Title != ""},
					Description:    pgtype.Text{String: myAd.Description, Valid: myAd.Description != ""},
					City:           pgtype.Text{String: myAd.City, Valid: myAd.City != ""},
					Neighborhood:   pgtype.Text{String: myAd.Neighborhood, Valid: myAd.Neighborhood != ""},
					HouseType:      sqlc.NullHouseType{HouseType: sqlc.HouseType(myAd.HouseType), Valid: myAd.HouseType != ""},
					Meterage:       pgtype.Int4{Int32: int32(myAd.Meterage), Valid: myAd.Meterage > 0},
					RoomsCount:     pgtype.Int4{Int32: int32(myAd.RoomsCount), Valid: myAd.RoomsCount > 0},
					Year:           pgtype.Int4{Int32: int32(myAd.Year), Valid: myAd.Year > 0},
					Floor:          pgtype.Int4{Int32: int32(myAd.Floor), Valid: myAd.Floor > 0},
					TotalFloors:    pgtype.Int4{Int32: int32(myAd.TotalFloors), Valid: myAd.TotalFloors > 0},
					HasWarehouse:   pgtype.Bool{Bool: myAd.HasWarehouse, Valid: true},
					HasElevator:    pgtype.Bool{Bool: myAd.HasElevator, Valid: true},
					Lat:            pgtype.Numeric{Int: big.NewInt(1), Valid: true},
					Lng:            pgtype.Numeric{Int: big.NewInt(2), Valid: true},
				}

				_, err := db.CreateAd(ctx, *ad)

				if err != nil {
					// add logger for here too
					fmt.Println(err)
					bot.Send(tgbotapi.NewMessage(userID, "Something went wrong"))

				} else {
					bot.Send(tgbotapi.NewMessage(userID, "House added successfully: \n\n"))
				}
			} else if userCtx.Command == "updatehouse" {
				// TODO: update ad
				bot.Send(tgbotapi.NewMessage(userID, "House updated successfully."))
			}
			delete(userContext, userID)

		} else {
			bot.Send(tgbotapi.NewMessage(userID, "Invalid value. Enter URL again"))
		}
	}
}
