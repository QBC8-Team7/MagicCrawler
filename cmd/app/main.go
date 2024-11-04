package main

import (
	"fmt"
	"github.com/QBC8-Team7/MagicCrawler/config"
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

	fmt.Println(cfg.ConnectionURI)

}
