package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/middleware"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service/auth"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service/resumes"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	err := service.JwtClient.InitJwtAuth()
	if err != nil {
		log.Fatal("Could not initialize jwt auth: ", err)
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		log.Fatal("Allowed origin for frontend not provided: ", err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(chiMiddleware.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprint("Test")))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			authStore := auth.NewStore(s.db)
			authHandler := auth.NewHandler(authStore)
			authHandler.AddRoutes(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.JwtAuthenticator())

			resumeStore := resumes.NewStore(s.db)
			resumeHandler := resumes.NewHandler(*resumeStore)
			resumeHandler.AddRoutes(r)
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
