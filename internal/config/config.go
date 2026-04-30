// Package config manages all application-level configurations.
package config

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Environment string

const (
	Development = "dev"
	Production  = "prod"
	Staging     = "staging"
)

// Config holds all application-level configurations.
type Config struct {
	DatabaseURL string `mapstructure:"KEVIN_DATABASE_URL"`
	RiotAPIKey  string `mapstructure:"KEVIN_RIOT_API_KEY"`
	Environment string `mapstructure:"KEVIN_ENV"`
	Port        int    `mapstructure:"KEVIN_PORT"`
}

// NewConfig loads configuration from environment variables.
func NewConfig() (*Config, error) {
	v := viper.New()

	if err := v.BindEnv("KEVIN_RIOT_API_KEY"); err != nil {
		return nil, fmt.Errorf("failed to bind KEVIN_RIOT_API_KEY: %w", err)
	}

	if err := v.BindEnv("KEVIN_DATABASE_URL"); err != nil {
		return nil, fmt.Errorf("failed to bind KEVIN_DATABASE_URL: %w", err)
	}

	if err := v.BindEnv("KEVIN_ENV"); err != nil {
		return nil, fmt.Errorf("failed to bind KEVIN_ENV: %w", err)
	}

	if err := v.BindEnv("KEVIN_PORT", "PORT", "KEVIN_PORT"); err != nil {
		return nil, fmt.Errorf("failed to bind PORT: %w", err)
	}

	config := Config{
		Port:        7331,
		Environment: Development,
	}
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal from viper: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	var errs []error

	if c.RiotAPIKey == "" {
		errs = append(errs, errors.New("KEVIN_RIOT_API_KEY is not set"))
	}

	if err := validateDatabaseURL(c.DatabaseURL); err != nil {
		errs = append(errs, fmt.Errorf("KEVIN_DATABASE_URL: %w", err))
	}

	if c.Environment != Development && c.Environment != Staging && c.Environment != Production {
		message := fmt.Sprintf("KEVIN_ENV must be either '%s' or '%s'", Development, Production)
		errs = append(errs, errors.New(message))
	}

	if 1024 > c.Port || c.Port > 65535 {
		errs = append(errs, errors.New("PORT must be between 1024 and 65535"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func validateDatabaseURL(raw string) error {
	_, err := pgx.ParseConfig(raw)
	if err != nil {
		return err
	}

	return nil
}
