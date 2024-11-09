package server

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v4"
)

func GenerateRoutes(s *BotServer) {
	s.Bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		data := c.Callback().Data
		fmt.Println(data)

		if data == "create_new_ad" {
			return s.Handler.HandleCategory(c)
		}
		// handle category selection
		if strings.Contains(data, "category_") {
			s.Handler.HandleCategorySelection(c)
		}
		if data == "ask_title" {
			s.Handler.HandleTitleAndDescription(c)
		}
		// if data == "watchlist_action" {
		// 	// return s.Handler.HandleWatchlist(c)
		// }
		return nil
	})
	// s.Bot.Handle("/start", s.Handler.HandleWelcome)
	// s.Bot.Handle("/hello", s.Handler.HandleHello)
	// s.Bot.Handle("/bye", s.Handler.HandleBye)
	// s.Bot.Handle(telebot.OnCallback, s.Handler.HandleCategoryCallback)

	s.Bot.Handle("/start", s.Handler.HandleStart)
	s.Bot.Handle("Create New Ad", s.Handler.HandleCategory)
	s.Bot.Handle("/category", s.Handler.HandleCategory)

	// s.Bot.Handle("/title", s.Handler.HandleTitleAndDescription)
	// s.Bot.Handle("/location", s.Handler.HandleLocation)
	// s.Bot.Handle("/neighborhood", s.Handler.HandleNeighborhood)
	// s.Bot.Handle("/propertydetails", s.Handler.HandlePropertyDetails)
	// s.Bot.Handle("/rooms", s.Handler.HandleRoomsCount)
	// s.Bot.Handle("/features", s.Handler.HandleFeatures)
	// s.Bot.Handle("/finalconfirmation", s.Handler.HandleFinalConfirmation)
	// s.Bot.Handle("/submit", s.Handler.HandleSubmitAd)

}
