package api_service

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/encoding/json"
	"github.com/storskegg/ampvoice/internal/dbModels"
	"gorm.io/gorm"
)

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

func getPart(ctx context.Context, db *gorm.DB) http.HandlerFunc {
	part, err := gorm.G[dbModels.Part](db).Where("id = ?", 1).First(ctx)
	if err != nil {
		log.Printf("Failed to find product: %s", err.Error())
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Product not found", http.StatusNotFound)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(part); err != nil {
			log.Error().Err(err).Msg("Failed to encode product")
			http.Error(w, "Failed to encode product", http.StatusInternalServerError)
		}
	}
}
