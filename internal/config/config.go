package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application-level configurations.
type Config struct {
	Environment        string `mapstructure:"KEVIN_ENVIRONMENT"`
	PostgresConnection string `mapstructure:"KEVIN_POSTGRES_CONNECTION"`
	RiotAPIKey         string `mapstructure:"KEVIN_RIOT_API_KEY"`
	Port               int    `mapstructure:"PORT"`
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("KEVIN_ENVIRONMENT", "development")
	v.SetDefault("PORT", 8080)

	_ = v.BindEnv("KEVIN_ENVIRONMENT")
	_ = v.BindEnv("KEVIN_POSTGRES_CONNECTION")
	_ = v.BindEnv("KEVIN_RIOT_API_KEY")
	_ = v.BindEnv("PORT")

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
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
