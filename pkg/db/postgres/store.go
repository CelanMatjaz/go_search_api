package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type PostgresStore struct {
	Db *sql.DB
}

func CreatePostgresStore(connectionString string) *PostgresStore {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Could not create postgres connection")
	}
	return &PostgresStore{Db: db}
}
