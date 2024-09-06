package api

import (
	"fmt"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service/auth"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service/resumes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	port string
	db   *db.DbConnection
}

func NewAPIServer(port string, db *db.DbConnection) *APIServer {
	return &APIServer{
		port: port,
		db:   db,
	}
}

func (s *APIServer) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint("Test")))
	})

	r.Route("/api/v1", func(r chi.Router) {
		authStore := auth.NewStore(s.db)
		authHandler := auth.NewHandler(authStore)
		authHandler.AddRoutes(r)

        resumeStore := resumes.NewStore(s.db)
        resumeHandler := resumes.NewHandler(*resumeStore)
        resumeHandler.AddRoutes(r)
	})

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", s.port),
		Handler: r,
	}

	fmt.Printf("Starting server on port %s\n", s.port)

	server.ListenAndServe()

	return nil
}
