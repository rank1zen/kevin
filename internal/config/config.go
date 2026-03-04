// config manages all application-level configurations. Only environment
// variables are supported. There are no default values. The two supported
// environments are development and production.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

const (
	development = "development"
	production  = "production"
)

// config is an unexported struct for unmarshaling from viper.
type config struct {
	KevinDatabaseURL string `mapstructure:"KEVIN_DATABASE_URL"`
	KevinRiotAPIKey  string `mapstructure:"KEVIN_RIOT_API_KEY"`
	KevinEnv         string `mapstructure:"KEVIN_ENV"`
	Port             int    `mapstructure:"PORT"`
}

// Config holds all application-level configurations.
type Config struct {
	kevinRiotAPIKey  string
	kevinDatabaseURL string
	kevinEnv         string
	port             int
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

	if err := v.BindEnv("PORT"); err != nil {
		return nil, fmt.Errorf("failed to bind PORT: %w", err)
	}

	output := config{}
	if err := v.Unmarshal(&output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal from viper: %w", err)
	}

	cfg := &Config{
		kevinRiotAPIKey:  output.KevinRiotAPIKey,
		kevinDatabaseURL: output.KevinDatabaseURL,
		kevinEnv:         output.KevinEnv,
		port:             output.Port,
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var errs []error

	if c.kevinRiotAPIKey == "" {
		errs = append(errs, errors.New("KEVIN_RIOT_API_KEY is not set"))
	}

	if err := validateDatabaseURL(c.kevinDatabaseURL); err != nil {
		errs = append(errs, fmt.Errorf("KEVIN_DATABASE_URL: %w", err))
	}

	if c.kevinEnv != "development" && c.kevinEnv != "production" {
		message := fmt.Sprintf("KEVIN_ENV must be either '%s' or '%s'", development, production)
		errs = append(errs, errors.New(message))
	}

	if 1024 > c.port || c.port > 65535 {
		errs = append(errs, errors.New("PORT must be between 1024 and 65535"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *Config) GetRiotAPIKey() string {
	return c.kevinRiotAPIKey
}

func (c *Config) GetDatabaseURL() string {
	return c.kevinDatabaseURL
}

func (c *Config) IsDevelopment() bool {
	return c.kevinEnv == development
}

func (c *Config) GetPort() int {
	return c.port
}

func validateDatabaseURL(raw string) error {
	s := strings.TrimSpace(raw)
	if s == "" {
		return errors.New("url is not set")
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	// (5) Scheme whitelist
	switch strings.ToLower(u.Scheme) {
	case "postgres", "postgresql":
	default:
		return fmt.Errorf("url must use postgres/postgresql scheme, got %q", u.Scheme)
	}

	// (4) Required URL parts
	if u.Host == "" {
		return errors.New("url must include a host")
	}
	// Path should be /dbname
	if u.Path == "" || u.Path == "/" {
		return errors.New("url must include a database name in path")
	}

	return nil
}
