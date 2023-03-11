// Package app is the setup package for all things API related. It calls properly initializes all other
// required API functions and starts the main API service.
package app

import (
	"github.com/clintjedwards/basecoat/internal/api"
	"github.com/clintjedwards/basecoat/internal/config"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/rs/zerolog/log"
)

// StartServices initializes all required services.
func StartServices(config *config.API) {
	if config.Development.UseLocalhostTLS || config.Development.BypassAuth {
		log.Warn().Msg("server in development mode; not for use in production")
	}

	newStorage, err := initStorage(config.Server)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init storage")
	}

	log.Info().Str("path", config.Server.StoragePath).Int("max_results_limit", config.Server.StorageResultsLimit).
		Msg("storage initialized")

	newAPI, err := api.NewAPI(config, newStorage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init api")
	}

	newAPI.StartAPIService()
}

// initStorage creates a storage object with the appropriate engine
func initStorage(config *config.Server) (storage.DB, error) {
	return storage.New(config.StoragePath, config.StorageResultsLimit)
}
