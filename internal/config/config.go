package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Port        int    `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
	Debug       bool   `mapstructure:"debug"`
}

// Load reads configuration from environment variables and validates required fields.
// CLI flags should be bound to Viper before calling this function.
func Load() (*Config, error) {
	viper.SetDefault("port", 8080)
	viper.SetDefault("debug", false)

	// Environment variable settings
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required configuration fields are set
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("database_url is required (set via --db flag or APP_DATABASE_URL env var)")
	}
	return nil
}
