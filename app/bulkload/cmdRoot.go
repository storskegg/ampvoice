package bulkload

import (
	"context"
	"errors"
	"fmt"
	"os"

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

var (
	paramFilePath string
)

func init() {
	cmdRoot.PersistentFlags().StringVar(&paramFilePath, "f", "", "Path to the JSON file to load")

	if err := cmdRoot.MarkPersistentFlagRequired("f"); err != nil {
		log.Fatal().Err(err).Msg("Failed to mark flag as required")
	}
}

var cmdRoot = &cobra.Command{
	Use:   "bulkload",
	Short: "Bulk load data into the database",
	Long:  "Bulk load data into the database from a JSON file",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	RunE: runRootE,
}

func runRootE(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	parts := new([]*dbModels.Part)

	data, err := os.ReadFile(paramFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err = json.Unmarshal(data, &parts); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	db := newDatabase()
	if err := dbModels.MigrateModels(db); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	if err = db.Transaction(func(tx *gorm.DB) error {
		for _, part := range *parts {
			if err := gorm.G[dbModels.Part](tx).Create(ctx, part); err != nil {
				return fmt.Errorf("failed to create part: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func newDatabase() *gorm.DB {
	log.Info().Msg("Connecting to database...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUsername, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connection established")

	log.Info().Msg("Migrating database...")
	// Migrate the schema
	err = dbModels.MigrateModels(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to migrate database")
	}
	log.Info().Msg("Database migrated")

	return db
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
