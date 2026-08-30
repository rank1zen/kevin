package runtime

import (
	"github.com/spf13/viper"
)

const (
	Development = "development"
	Production  = "production"
	Staging     = "staging"
)

// Config holds all application-level configurations.
type Config struct {
	// DatabaseURL is the URL of the database.
	DatabaseURL string `mapstructure:"KEVIN_DATABASE_URL"`

	// RiotAPIKey is the API key for the Riot Games API.
	RiotAPIKey string `mapstructure:"KEVIN_RIOT_API_KEY"`

	// Environment is the environment in which the application is running. Defaults to Production.
	Environment string `mapstructure:"KEVIN_ENV"`

	// Port is the port on which the server will listen. Defaults to 7331.
	Port int `mapstructure:"KEVIN_PORT"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	err := v.BindEnv("KEVIN_RIOT_API_KEY")
	if err != nil {
		return nil, err
	}

	err = v.BindEnv("KEVIN_DATABASE_URL")
	if err != nil {
		return nil, err
	}

	err = v.BindEnv("KEVIN_ENV")
	if err != nil {
		return nil, err
	}

	err = v.BindEnv("KEVIN_PORT", "PORT", "KEVIN_PORT")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Environment: Production,
		Port:        7331,
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return mergeToDefaultConfig(cfg), nil
}

func mergeToDefaultConfig(cfg *Config) *Config {
	var environment = Production
	if cfg.Environment == Production || cfg.Environment == Development || cfg.Environment == Staging {
		environment = cfg.Environment
	}

	var port = 7331
	if cfg.Port != 0 {
		port = cfg.Port
	}

	return &Config{
		DatabaseURL: cfg.DatabaseURL,
		RiotAPIKey:  cfg.RiotAPIKey,
		Environment: environment,
		Port:        port,
	}
}
