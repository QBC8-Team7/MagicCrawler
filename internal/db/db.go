package db

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"sync"
)

var conn *sql.DB
var once sync.Once

// GetDBConnection gets the uri of db and return the connection; it creates the connection only once
func GetDBConnection(uri string) *sql.DB {
	once.Do(func() {
		c, err := sql.Open("postgres", uri)
		if err != nil {
			log.Fatalf("Failed to connect to db: %v", err)
		}

		if err = conn.Ping(); err != nil {
			log.Fatalf("Failed to ping db: %v", err)
		}
		conn = c
	})

	return conn
}
