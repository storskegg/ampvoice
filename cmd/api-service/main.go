package main

import (
	"github.com/rs/zerolog/log"
	apiservice "github.com/storskegg/ampvoice/app/api-service"
)

func main() {
	if err := apiservice.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run API service")
	}
}
