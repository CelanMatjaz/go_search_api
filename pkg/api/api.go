package api

import (
	"fmt"
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/handlers"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	port  string
	store db.Store
}

func CreateServer(port string, store db.Store) *Server {
	return &Server{
		port:  port,
		store: store,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()

	r.Use(chiMiddleware.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		authHandler := handlers.NewAuthHandler(s.store)
		authHandler.AddRoutes(r)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator(s.store))

			tagHandler := handlers.NewTagHandler(s.store)
			tagHandler.AddRoutes(r)
		})
	})

	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", s.port),
		Handler: r,
	}

	fmt.Printf("Starting server on port %s\n", s.port)

	server.ListenAndServe()

	return nil
}
