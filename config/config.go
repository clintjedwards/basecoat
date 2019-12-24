package config

import (
	"github.com/kelseyhightower/envconfig"
)

// DatabaseConfig refers to database connection settings
type DatabaseConfig struct {
	// The database engine used by the backend
	// possible values are: googleDatastore
	Engine          string `envconfig:"database_engine" default:"googleDatastore"`
	GoogleDatastore *GoogleDatastoreConfig
}

// GoogleDatastoreConfig represents google firebase datastore configuration
// https://cloud.google.com/datastore/docs/concepts/overview
type GoogleDatastoreConfig struct {
	ProjectID string `envconfig:"database_google_datastore_project_id" default:"test"`
	// Use local emulator as a test database
	EmulatorHost string `envconfig:"database_google_datastore_emulator_host" default:"localhost:8000"`
	// Timeout for RPC calls to the database. Will accept duration string as noted in
	// https://golang.org/pkg/time/#ParseDuration
	Timeout string `envconfig:"database_google_datastore_timeout" default:"10s"`
}

// Config refers to general application configuration
type Config struct {
	// Debug is useful for development builds; turns on extra/better logging
	Debug       bool   `envconfig:"debug" default:"false"`
	TLSCertPath string `envconfig:"tls_cert_path" default:"./localhost.crt"`
	TLSKeyPath  string `envconfig:"tls_key_path" default:"./localhost.key"`
	URL         string `envconfig:"url" default:"localhost:8080"`
	Frontend    *FrontendConfig
	Backend     *BackendConfig
	Database    *DatabaseConfig
	CommandLine *CommandLineConfig
	Metrics     *MetricsConfig
}

// FrontendConfig represents configuration for frontend basecoat
type FrontendConfig struct {
	Enable bool `envconfig:"frontend_enable" default:"true"`
}

// BackendConfig represents configuration for backend basecoat grpc service
type BackendConfig struct {
	IDLength  int    `envconfig:"backend_id_length" default:"5"`          // the length of all randomly generated ids
	SecretKey string `envconfig:"backend_secret_key" default:"testtoken"` // secret key used to encrypt api tokens
}

// MetricsConfig represents configuration for the metrics endpoint
type MetricsConfig struct {
	Endpoint string `envconfig:"metrics_endpoint" default:"localhost:8082"`
}

// CommandLineConfig represents configuration for cli application
type CommandLineConfig struct {
	Token string `envconfig:"token" default:""`
}

// FromEnv pulls configuration from environment variables
func FromEnv() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
