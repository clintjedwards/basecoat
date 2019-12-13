package utils

import (
	"log"

	"github.com/clintjedwards/basecoat/config"
	"go.uber.org/zap"
)

// Log returns a configured logger
func Log() *zap.SugaredLogger {
	// The call to FromEnv here doubles the runtime of this
	// method (~80ns to ~144ns).
	// Can be removed for performance, but is here for dev convenience
	config, err := config.FromEnv()
	if err != nil {
		log.Printf("could not get config: %v", err)
	}

	var logger *zap.Logger
	if config.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Printf("could not init logger: %v", err)
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	return sugar
}
