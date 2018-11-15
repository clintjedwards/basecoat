package config

import (
	"github.com/kelseyhightower/envconfig"
)

//Config refers to command line application configuration
type Config struct {
	Debug       bool   `envconfig:"debug" default:"false"`
	BasecoatURL string `envconfig:"basecoat_url" default:"http://localhost:8080"`
}

//FromEnv pulls configuration from environment variables
func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
