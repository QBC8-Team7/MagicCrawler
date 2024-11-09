package server

import (
	"strconv"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandHandler interface {
	HandleHello(update tgbotapi.Update) error
	HandleBye(update tgbotapi.Update) error
	HandleWatchlist(update tgbotapi.Update) error
}

type Handlers struct {
	Logger *logger.AppLogger
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
	Lat            float64
	Lng            float64
}

func (h *Handlers) StartFlow(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Welcome! Let's create your home ad. Please choose the category.")
	categories := []string{"Apartment", "House", "Villa", "Studio"}

	var inlineButtons []tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		inlineButtons = append(inlineButtons, tgbotapi.NewInlineKeyboardButtonData(category, "category_"+category))
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(inlineButtons)

	msg.ReplyMarkup = menu

	bot.Send(msg)
}

func (h *Handlers) HandleCategorySelection(bot *tgbotapi.BotAPI, chatID int64, category string) {
	h.Ad.Category = category

	msg := tgbotapi.NewMessage(chatID, "Please enter the title of your ad (e.g., 2-bedroom apartment for rent):")
	bot.Send(msg)
}

func (h *Handlers) HandleTitle(bot *tgbotapi.BotAPI, chatID int64, title string) {
	h.Ad.Title = title

	msg := tgbotapi.NewMessage(chatID, "Please provide a description for your ad:")
	bot.Send(msg)
}

func (h *Handlers) HandleDescription(bot *tgbotapi.BotAPI, chatID int64, description string) {
	h.Ad.Description = description

	msg := tgbotapi.NewMessage(chatID, "Please enter the city where the property is located:")
	bot.Send(msg)
}

func (h *Handlers) HandleCity(bot *tgbotapi.BotAPI, chatID int64, city string) {
	h.Ad.City = city

	msg := tgbotapi.NewMessage(chatID, "Please enter the neighborhood:")
	bot.Send(msg)
}

func (h *Handlers) HandleNeighborhood(bot *tgbotapi.BotAPI, chatID int64, neighborhood string) {
	h.Ad.Neighborhood = neighborhood

	msg := tgbotapi.NewMessage(chatID, "Please enter the floor number of the property:")
	bot.Send(msg)
}

func (h *Handlers) HandleFloor(bot *tgbotapi.BotAPI, chatID int64, floor string) {
	floorInt, err := strconv.Atoi(floor)
	if err == nil {
		h.Ad.Floor = floorInt
	}

	msg := tgbotapi.NewMessage(chatID, "Please enter the total number of floors in the building:")
	bot.Send(msg)
}

func (h *Handlers) HandleTotalFloors(bot *tgbotapi.BotAPI, chatID int64, totalFloors string) {
	totalFloorsInt, err := strconv.Atoi(totalFloors)
	if err == nil {
		h.Ad.TotalFloors = totalFloorsInt
	}

	msg := tgbotapi.NewMessage(chatID, "Does the property have a warehouse? (Yes/No)")
	bot.Send(msg)
}

func (h *Handlers) HandleBooleanInput(bot *tgbotapi.BotAPI, chatID int64, response string, field *bool) {
	if response == "Yes" {
		*field = true
	} else {
		*field = false
	}

	msg := tgbotapi.NewMessage(chatID, "Does the property have an elevator? (Yes/No)")
	bot.Send(msg)
}
