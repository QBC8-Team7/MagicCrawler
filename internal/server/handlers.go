package server

import (
	"fmt"

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
}

func (h *Handlers) HandleStart(m *tgbotapi.Message, c *tgbotapi.BotAPI) error {
	startMessage := "Welcome! 😊 Would you like to create a new ad? Click below to begin."

	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("create ad", "create_ad"),
		),
	)
	msg := tgbotapi.NewMessage(m.Chat.ID, startMessage)
	msg.ReplyMarkup = numericKeyboard

	c.Send(msg)
	return nil
}

func (h *Handlers) HandleCategory(update tgbotapi.Update) error {
	categories := []string{"Apartment", "House", "Villa", "Studio"}

	var buttons []tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(category, "category_"+category))
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Please select an ad category:")
	editMsg.ReplyMarkup = &menu
	_, err := h.Bot.Send(editMsg)
	return err
}

func (h *Handlers) HandleCategorySelection(update tgbotapi.Update) error {
	selectedCategory := update.CallbackQuery.Data

	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Next: Title", "ask_title"),
		),
	)
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, fmt.Sprintf("You selected: %s", selectedCategory))
	editMsg.ReplyMarkup = &menu
	_, err := h.Bot.Send(editMsg)
	return err
}

func (h *Handlers) HandleTitleAndDescription(update tgbotapi.Update) error {
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Please enter the title of the ad (e.g., 2-bedroom apartment for rent):")
	_, err := h.Bot.Send(editMsg)
	return err
}

// Additional functions for handling other commands (HandleLocation, HandleNeighborhood, etc.) would follow a similar pattern,
// using tgbotapi.NewMessage or tgbotapi.NewEditMessageText depending on the interaction type.
