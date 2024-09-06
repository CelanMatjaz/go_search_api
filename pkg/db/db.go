package db

import (
	"database/sql"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/assert"
	_ "github.com/lib/pq"
)

type DbConnection struct {
	DB *sql.DB
}

func NewDbConnection(connectionString string) *DbConnection {
	db, err := sql.Open("postgres", connectionString)
	assert.AssertError(err, "Could not connect to DB")

	return &DbConnection{DB: db}
}
