package main

import (
	"database/sql"
	"flag"
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
	envFile := flag.String("env", "local.env", "env file")
	flag.Parse()

	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Fatal("Provide a subcommand, either of 'up', 'down' or 'reset'")
	}

	flag.Parse()

	godotenv.Load(*envFile)
	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	path, _ := filepath.Abs("migrations")
	migration, err := migrate.NewWithDatabaseInstance(fmt.Sprint("file://", path), "postgres", driver)
	if err != nil {
		panic(err)
	}

	switch flag.Args()[0] {
	case "up":
		upCmd.Parse(os.Args[2:])
		if err := migration.Up(); err != nil {
			log.Fatalf("Could not migrate up, %s", err.Error())
		}
		break

	case "down":
		downCmd.Parse(os.Args[2:])
		if err := migration.Down(); err != nil {
			log.Fatalf("Could not migrate down, %s", err.Error())
		}
		break

	case "reset":
		resetCmd.Parse(os.Args[2:])
		if err := migration.Drop(); err != nil {
			log.Fatalf("Could not drop db, %s", err.Error())
		}
		break

	default:
		log.Fatalf("Only subcommands 'up', 'down' and 'reset' supported, command '%s' provided", flag.Args()[0])
	}
}
