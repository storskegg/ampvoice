package api_service

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/storskegg/ampvoice/internal/dbModels"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
