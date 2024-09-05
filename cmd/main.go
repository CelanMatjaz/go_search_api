package main

import (
	"os"

	"github.com/CelanMatjaz/job_application_tracker_api/cmd/api"
	"github.com/CelanMatjaz/job_application_tracker_api/cmd/db"
)

func main()  {
	db := db.NewDbConnection(os.Getenv("CONNECTION_STRING"))

	server := api.NewAPIServer("8080", db)
	server.Start()
}

