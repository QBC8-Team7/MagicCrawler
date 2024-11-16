package main

import (
	"context"
	"flag"
	"fmt"
	myredis "github.com/QBC8-Team7/MagicCrawler/pkg/redis"
	"log"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/jackc/pgx/v5"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/server"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db"
)

func main() {
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config error: %w", err))
	}

	dbContext := context.Background()

	redisClient := myredis.NewRedisClient(dbContext, conf.Redis.Host, conf.Redis.Port, conf.Redis.Password, conf.Redis.DB)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("Failed to close Redis client: %v", err)
		}
	}()

	dbUri := db.GetDbUri(conf)
	dbConn, err := db.GetDBConnection(dbContext, dbUri)
	if err != nil {
		log.Fatalln(fmt.Errorf("could not connect to database: %w", err))
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatalln(fmt.Errorf("could not close connection with database: %w", err))
		}
	}(dbConn, dbContext)

	dbQueries := sqlc.New(dbConn)

	s, err := server.NewServer(dbContext, conf, dbQueries, redisClient)
	if err != nil {
		log.Fatal(fmt.Errorf("could not start server: %w", err))
	}

	err = s.Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("error while running server: %w", err))
	}
}
