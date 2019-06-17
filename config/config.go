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
	Debug       bool   `envconfig:"debug" default:"false"`
	TLSCertPath string `envconfig:"tls_cert_path" default:"./localhost.crt"`
	TLSKeyPath  string `envconfig:"tls_key_path" default:"./localhost.key"`
	Frontend    *FrontendConfig
	Backend     *BackendConfig
	Database    *DatabaseConfig
	CommandLine *CommandLineConfig
}

// FrontendConfig represents configuration for frontend basecoat
type FrontendConfig struct {
	Enable bool `envconfig:"frontend_enable" default:"true"`
	// This envvar is not used from this config but is here for completeness
	// it is set in the makefile and is injected into the js code at build time.
	// It controls where the frontend client should look for the gprc backend
	APIHost string `envconfig:"frontend_api_host" default:"https://localhost:8080"`
}

// BackendConfig represents configuration for backend basecoat grpc service
type BackendConfig struct {
	IDLength  int    `envconfig:"backend_id_length" default:"5"`          // the length of all randomly generated ids
	SecretKey string `envconfig:"backend_secret_key" default:"testtoken"` // secret key used to encrypt api tokens
	HTTPURL   string `envconfig:"backend_http_url" default:"localhost:8080"`
	GRPCURL   string `envconfig:"backend_grpc_url" default:"localhost:8081"`
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
