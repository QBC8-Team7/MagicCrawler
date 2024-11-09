package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/server"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", err)
	}

	db_uri := db.GetDbUri(cfg)
	db, err := db.GetDBConnection(db_uri, cfg.PgDriver)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	s := server.NewServer(cfg)

	go func() {
		fmt.Println("Bot Server Started...")
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down bot...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Bot.StopReceivingUpdates()

	<-ctx.Done()

	s.Logger.Info("Bot exited gracefully.")

}
