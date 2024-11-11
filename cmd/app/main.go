package main

import (
	"context"
	"log"

	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/jackc/pgx/v5"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/server"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("could not read config file: ", err)
	}

	ctx := context.Background()

	dbUri := db.GetDbUri(cfg)
	dbConn, err := db.GetDBConnection(ctx, dbUri)
	if err != nil {
		log.Fatal("could not connect to db: ", err)
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Fatal("could not close connection:", err)
		}
	}(dbConn, ctx)

	dbQueries := sqlc.New(dbConn)

	s := server.NewServer(ctx, cfg, dbQueries)

	s.Run()

}
