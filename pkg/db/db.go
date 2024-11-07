package db

import (
	"fmt"
	"sync"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	conn *sqlx.DB
	once sync.Once
)

func GetDBConnection(uri, driver_name string) (*sqlx.DB, error) {
	var initErr error
	once.Do(func() {
		db, err := sqlx.Connect(driver_name, uri)

		if err != nil {
			initErr = fmt.Errorf("failed to connect to db: %v", err)
			return
		}
		conn = db
	})

	if initErr != nil {
		return nil, initErr
	}

	return conn, nil
}

func GetDbUri(cfg *config.Config) string {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlPassword,
	)

	return dataSourceName

}
