package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Provide a subcommand, either of 'up', 'down' or 'reset'")
	}

	godotenv.Load("dev.env")
	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	path, _ := filepath.Abs("migrations")
	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprint("file://", path), "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "up":
		if err := migration.Up(); err != nil {
			panic(err)
		}
		break

	case "down":
		if err := migration.Down(); err != nil {
			panic(err)
		}
		break

	case "reset":
		if err := migration.Drop(); err != nil {
			panic(err)
		}
		break

	default:
		log.Fatal("Only subcommands 'up', 'down' and 'reset' supported")
	}
}
