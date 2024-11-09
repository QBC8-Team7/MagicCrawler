package server

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Example handler for inline button callback queries
func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, h *Handlers) {
	callback := update.CallbackQuery
	// Handle different callback data values here
	fmt.Println("=============", callback.Data)
	switch callback.Data {
	case "create_ad":
		fmt.Println("creaaaaaaaaaate ad")
		// msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Your watchlist is empty!")
		// if _, err := bot.Send(msg); err != nil {
		// 	log.Printf("Failed to send callback message: %v", err)
		// }
		// Send an answer to the callback query to remove the "waiting circle"
		// bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Showing your watchlist"))
	default:
		// Handle other callbacks if needed
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Unknown action"))
	}
}

// func GenerateRoutes(s *BotServer) {
// 	s.Bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
// 		data := c.Callback().Data
// 		fmt.Println(data)

// 		if data == "create_new_ad" {
// 			return s.Handler.HandleCategory(c)
// 		}
// 		// handle category selection
// 		if strings.Contains(data, "category_") {
// 			s.Handler.HandleCategorySelection(c)
// 		}
// 		if data == "ask_title" {
// 			s.Handler.HandleTitleAndDescription(c)
// 		}
// 		// if data == "watchlist_action" {
// 		// 	// return s.Handler.HandleWatchlist(c)
// 		// }
// 		return nil
// 	})
// 	// s.Bot.Handle("/start", s.Handler.HandleWelcome)
// 	// s.Bot.Handle("/hello", s.Handler.HandleHello)
// 	// s.Bot.Handle("/bye", s.Handler.HandleBye)
// 	// s.Bot.Handle(telebot.OnCallback, s.Handler.HandleCategoryCallback)

// 	s.Bot.Handle("/start", s.Handler.HandleStart)
// 	s.Bot.Handle("Create New Ad", s.Handler.HandleCategory)
// 	s.Bot.Handle("/category", s.Handler.HandleCategory)

// 	// s.Bot.Handle("/title", s.Handler.HandleTitleAndDescription)
// 	// s.Bot.Handle("/location", s.Handler.HandleLocation)
// 	// s.Bot.Handle("/neighborhood", s.Handler.HandleNeighborhood)
// 	// s.Bot.Handle("/propertydetails", s.Handler.HandlePropertyDetails)
// 	// s.Bot.Handle("/rooms", s.Handler.HandleRoomsCount)
// 	// s.Bot.Handle("/features", s.Handler.HandleFeatures)
// 	// s.Bot.Handle("/finalconfirmation", s.Handler.HandleFinalConfirmation)
// 	// s.Bot.Handle("/submit", s.Handler.HandleSubmitAd)

// }
