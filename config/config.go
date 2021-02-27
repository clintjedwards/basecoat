package config

import (
	"github.com/kelseyhightower/envconfig"
)

// DatabaseConfig refers to database connection settings
type DatabaseConfig struct {
	IDLength int    `envconfig:"database_id_length" default:"5"` // the length of all randomly generated ids
	Path     string `envconfig:"database_path" default:"/tmp/basecoat.db"`
}

// Config refers to general application configuration
type Config struct {
	Debug       bool   `envconfig:"debug" default:"false"`
	LogLevel    string `envconfig:"loglevel" default:"info"`
	TLSCertPath string `envconfig:"tls_cert_path" default:"./localhost.crt"` // does not apply if certmagic is enabled
	TLSKeyPath  string `envconfig:"tls_key_path" default:"./localhost.key"`  // does not apply if certmagic is enabled
	URL         string `envconfig:"url" default:"localhost:8080"`            // does not apply if certmagic is enabled
	CertMagic   *CertMagicConfig
	Frontend    *FrontendConfig
	Backend     *BackendConfig
	Database    *DatabaseConfig
	CommandLine *CommandLineConfig
	Metrics     *MetricsConfig
}

// CertMagicConfig allows the automation of generating tls certificates
type CertMagicConfig struct {
	// this enables cert magic to perform https cert renewal and redirection for the application
	Enable bool   `envconfig:"certmagic_enable" default:"false"`
	Email  string `envconfig:"certmagic_email"`
	Domain string `envconfig:"certmagic_domain" default:"basecoat.clintjedwards.com"`
}

// FrontendConfig represents configuration for frontend basecoat
type FrontendConfig struct {
	Enable bool `envconfig:"frontend_enable" default:"true"`
}

// BackendConfig represents configuration for backend basecoat grpc service
type BackendConfig struct {
	// duration limit on user requested api token, after limit token will expire
	APITokenDurationLimit  int64  `envconfig:"backend_api_token_duration_limit" default:"946708560"` // 946708560 = 30 years
	AdminToken             string `envconfig:"backend_admin_token" default:"admin"`                  // used to allow admin functions
	SearchIndexRebuildTime int64  `envconfig:"backend_search_index_rebuild_time" default:"600"`      // how often the search index rebuilds; in seconds
	SecretKey              string `envconfig:"backend_secret_key" default:"testtoken"`               // secret key used to encrypt api tokens
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
