//go:generate go run frontend/generate.go

// Basecoat is an internal formula management system
package main

import (
	"log"

	"github.com/clintjedwards/basecoat/cmd"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/toolkit/logger"
)

func init() {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	logger.InitGlobalLogger(config.LogLevel, config.Debug)
}

func main() {
	cmd.RootCmd.Execute()
}
