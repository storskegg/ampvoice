package api_service

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	errMissingEnvironmentVariable = errors.New("missing environment variable")
)

var (
	mysqlUsername = ""
	mysqlPassword = ""
	mysqlDatabase = ""
	mysqlHost     = ""
	mysqlPort     = ""
)

var cmdRoot = cobra.Command{
	Use:   "api-service",
	Short: "API service",
	Long:  "API service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	RunE: runRootE,
}

func Run() error {
	return cmdRoot.Execute()
}

func runRootE(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting server...")

	ctx := context.Background()

	db := newDatabase()

	log.Info().Msg("Starting server...")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Get("/parts", getParts(ctx, db))
	})

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(r))
	// start server
	if err := http.ListenAndServe(":8000", enhancedRouter); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	return nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func initializeConfig(cmd *cobra.Command) error {
	var ok bool
	mysqlUsername, ok = os.LookupEnv("MYSQL_USER")
	if !ok || mysqlUsername == "" {
		return errors.Join(errMissingEnvironmentVariable, errors.New("MYSQL_USER"))
		//log.Fatal().Err(errMissingEnvironmentVariable).Msg("MYSQL_USER")
	}
	mysqlPassword, ok = os.LookupEnv("MYSQL_PASSWORD")
	if !ok || mysqlPassword == "" {
		return errors.Join(errMissingEnvironmentVariable, errors.New("MYSQL_PASSWORD"))
		//log.Fatal().Err(errMissingEnvironmentVariable).Msg("MYSQL_PASSWORD")
	}
	mysqlDatabase, ok = os.LookupEnv("MYSQL_DATABASE")
	if !ok || mysqlDatabase == "" {
		return errors.Join(errMissingEnvironmentVariable, errors.New("MYSQL_DATABASE"))
		//log.Fatal().Err(errMissingEnvironmentVariable).Msg("MYSQL_DATABASE")
	}
	mysqlHost, ok = os.LookupEnv("MYSQL_HOST")
	if !ok || mysqlHost == "" {
		return errors.Join(errMissingEnvironmentVariable, errors.New("MYSQL_HOST"))
	}
	mysqlPort, ok = os.LookupEnv("MYSQL_PORT")
	if !ok || mysqlPort == "" {
		return errors.Join(errMissingEnvironmentVariable, errors.New("MYSQL_PORT"))
	}

	return nil
}
