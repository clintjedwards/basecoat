package config

import (
	"github.com/kelseyhightower/envconfig"
)

//DatabaseConfig refers to database connection settings
type DatabaseConfig struct {
	Name     string `envconfig:"database_name" default:"basecoat"`
	URL      string `envconfig:"database_url" default:"localhost:5432"`
	User     string `envconfig:"database_user" default:"basecoat"`
	Password string `envconfig:"database_password" default:"basecoat"`
}

//Config refers to general application configuration
type Config struct {
	ServerURL string `envconfig:"server_url" default:"localhost:8080"`
	SecretKey string `envconfig:"secret_key" default:"testtoken"` //Used as single login for website
	Debug     bool   `envconfig:"debug" default:"false"`
	Frontend  bool   `envconfig:"frontend" default:"true"`
	Database  *DatabaseConfig
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
