package main

import (
	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/app"
	"log"
)

func main() {
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatal("Could not parse config file: ", err)
	}

	s := app.NewServer(cfg.Token)

	s.Bot.Start()
}
