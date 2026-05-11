package api_service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/cobra"
	"github.com/storskegg/ampvoice/internal/dbModels"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeConfig(cmd)
	},
	RunE: runRootE,
}

func Run() error {
	return cmdRoot.Execute()
}

func runRootE(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting server...")

	log.Info().Msg("Connecting to database...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connection established")

	ctx := context.Background()

	log.Info().Msg("Migrating database...")
	// Migrate the schema
	err = dbModels.MigrateModels(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to migrate database")
	}
	log.Info().Msg("Database migrated")

	log.Info().Msg("Starting server...")
	router := mux.NewRouter()
	router.HandleFunc("/api/parts", getParts(ctx, db)).Methods("GET")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
	// start server
	err = http.ListenAndServe(":8000", enhancedRouter)
	if err != nil {
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

func getParts(ctx context.Context, db *gorm.DB) http.HandlerFunc {
	parts, err := gorm.G[dbModels.Part](db).Find(ctx) // find product with integer primary key
	if err != nil {
		log.Printf("Failed to find product: %s", err.Error())
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Product not found", http.StatusNotFound)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(parts); err != nil {
			log.Error().Err(err).Msg("Failed to encode product")
			http.Error(w, "Failed to encode product", http.StatusInternalServerError)
		}
	}
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
