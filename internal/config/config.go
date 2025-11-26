package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application-level configurations.
type Config struct {
	Environment        string `mapstructure:"ENVIRONMENT"`
	PostgresConnection string `mapstructure:"POSTGRES_CONNECTION"`
	RiotAPIKey         string `mapstructure:"RIOT_API_KEY"`
	Address            string `mapstructure:"ADDRESS"`
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("ADDRESS", "0.0.0.0:4001")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	viper.SetEnvPrefix("KEVIN")
	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// validateConfig checks if required configuration fields are set.
func validateConfig(cfg *Config) error {
	if cfg.RiotAPIKey == "" {
		return errors.New("RIOT_API_KEY is not set")
	}
	if cfg.PostgresConnection == "" {
		return errors.New("POSTGRES_CONNECTION is not set")
	}
	return nil
}

// IsDevelopment checks if the current environment is development.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction checks if the current environment is production.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
