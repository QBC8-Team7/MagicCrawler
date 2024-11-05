package main

import (
	"log"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	s := app.NewServer(cfg.Token)
	s.Serve()
}
