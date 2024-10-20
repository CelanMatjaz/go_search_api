package main

import (
	"log"
	"os"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/api"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db/postgres"
)

func main() {
	connectionString := os.Getenv("CONNECTION_STRING")
	port := os.Getenv("PORT")

	store := postgres.NewPostgresStore(connectionString)
	server := api.CreateServer(port, store)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
