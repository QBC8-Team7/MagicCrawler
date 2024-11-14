package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"

	"github.com/QBC8-Team7/MagicCrawler/config"
)

var (
	conn *pgx.Conn
	once sync.Once
)

func GetDBConnection(ctx context.Context, uri string) (c *pgx.Conn, e error) {
	once.Do(func() {
		pgxConn, err := pgx.Connect(ctx, uri)
		if err != nil {
			c, e = nil, err
			return
		}
		conn = pgxConn
	})

	return conn, e
}

func GetDbUri(cfg *config.Config) string {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Dbname,
		cfg.Postgres.Password,
	)

	return dataSourceName
}
