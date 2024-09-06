package main

import (
	"os"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/api"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
)

func main()  {
	database := db.NewDbConnection(os.Getenv("CONNECTION_STRING"))

	server := api.NewAPIServer("8080", database)
	server.Start()
}

