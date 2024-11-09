package server

import (
	"fmt"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
	"gopkg.in/telebot.v4"
)

type CommandHandler interface {
	HandleHello(c telebot.Context) error
	HandleBye(c telebot.Context) error
	HandleWatchlist(c telebot.Context) error
}

type Handlers struct {
	Logger *logger.AppLogger
}

var (
	// Inline menu for inline buttons
	inlineMenu = &telebot.ReplyMarkup{}

	// Define inline button with callback data
	btnWatchlist = inlineMenu.Data("📋 My Watchlist", "watchlist_action")
)

// Initialize inline menu layout with the button
func init() {
	inlineMenu.Inline(
		inlineMenu.Row(btnWatchlist),
	)
}
func (h *Handlers) HandleStart(c telebot.Context) error {
	startMessage := "Welcome! 😊 Would you like to create a new ad? Click below to begin."

	menu := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				telebot.InlineButton{
					Text: "Create New Ad",
					Data: "create_new_ad", // This is the callback data, you can check it in your handler
				},
			},
		},
	}

	// Send the welcome message with the inline menu
	return c.Send(startMessage, menu)
}

func (h *Handlers) HandleCategory(c telebot.Context) error {
	categories := []string{"Apartment", "House", "Villa", "Studio"}

	menu := &telebot.ReplyMarkup{InlineKeyboard: make([][]telebot.InlineButton, 0)}

	var row []telebot.InlineButton
	for _, category := range categories {
		row = append(row, telebot.InlineButton{
			Text: category,
			Data: "category_" + category,
		})
	}
	menu.InlineKeyboard = append(menu.InlineKeyboard, row)

	return c.Edit("Please select an ad category:", menu)
}

func (h *Handlers) HandleCategorySelection(c telebot.Context) error {
	selectedCategory := c.Callback().Data

	menu := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				telebot.InlineButton{
					Text: "Next: Title",
					Data: "ask_title",
				},
			},
		},
	}
	return c.Edit(fmt.Sprintf("You selected: %s", selectedCategory), menu)
}

func (h *Handlers) HandleTitleAndDescription(c telebot.Context) error {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
	// btnNext := menu.Text("Next")
	// menu.Reply(menu.Row(btnNext))

	// Ask for the title
	return c.Edit("Please enter the title of the ad (e.g., 2-bedroom apartment for rent):", menu)
}

// func (h *Handlers) HandleLocation(c telebot.Context) error {
// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnNext := menu.Text("Next")
// 	menu.Reply(menu.Row(btnNext))

// 	// Ask for the city
// 	return c.Send("Please enter the city where the property is located:", menu)
// }

// func (h *Handlers) HandleNeighborhood(c telebot.Context) error {
// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnNext := menu.Text("Next")
// 	menu.Reply(menu.Row(btnNext))

// 	// Ask for the neighborhood
// 	return c.Send("Please enter the neighborhood:", menu)
// }

// func (h *Handlers) HandlePropertyDetails(c telebot.Context) error {
// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnNext := menu.Text("Next")
// 	menu.Reply(menu.Row(btnNext))

// 	// Ask for the meterage
// 	return c.Send("Please enter the size of the property in square meters:", menu)
// }

// func (h *Handlers) HandleRoomsCount(c telebot.Context) error {
// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnNext := menu.Text("Next")
// 	menu.Reply(menu.Row(btnNext))

// 	// Ask for the number of rooms
// 	return c.Send("How many rooms does the property have?", menu)
// }

// func (h *Handlers) HandleFeatures(c telebot.Context) error {
// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnYes := menu.Text("Yes")
// 	btnNo := menu.Text("No")

// 	menu.Reply(menu.Row(btnYes, btnNo))

// 	return c.Send("Does the property have a warehouse?", menu)
// }

// func (h *Handlers) HandleFinalConfirmation(c telebot.Context) error {
// 	summary := "Please review your ad details:\n\n" +
// 		"Category: Apartment\n" +
// 		"Title: Beautiful 2-bedroom apartment\n" +
// 		"Description: A spacious apartment with modern amenities...\n" +
// 		"Location: Tehran, District 5\n" +
// 		"Size: 120 m²\n" +
// 		"Rooms: 2\n" +
// 		"Has Elevator: Yes\n" +
// 		"Has Warehouse: No"

// 	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
// 	btnConfirm := menu.Text("Confirm")
// 	btnCancel := menu.Text("Cancel")

// 	menu.Reply(menu.Row(btnConfirm, btnCancel))

// 	return c.Send(summary+"\n\nDo you want to submit this ad?", menu)
// }

// func (h *Handlers) HandleSubmitAd(c telebot.Context) error {
// 	// Here you would typically save the ad data to a database or API
// 	return c.Send("Your ad has been submitted successfully! 🎉", nil)
// }
