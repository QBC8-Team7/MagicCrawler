package server

func GenerateRoutes(s *BotServer) {
	s.Bot.Handle("/hello", s.Handler.HandleHello)
	s.Bot.Handle("/bye", s.Handler.HandleBye)
}
