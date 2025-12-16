package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application-level configurations.
type Config struct {
	kevinPostgresConnection string
	kevinRiotAPIKey         string
	env                     string
	port                    int
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("ENV", "prod")
	v.SetDefault("PORT", 8080)

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	_ = v.BindEnv("KEVIN_POSTGRES_CONNECTION")
	_ = v.BindEnv("KEVIN_RIOT_API_KEY")
	_ = v.BindEnv("ENV")
	_ = v.BindEnv("PORT")

	output := config{}
	if err := v.Unmarshal(&output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg := output.ToConfig()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// IsDevelopment checks if the current environment is development.
func (c *Config) IsDevelopment() bool {
	return c.env == "dev"
}

func (c *Config) Validate() error {
	if c.kevinRiotAPIKey == "" {
		return errors.New("KEVIN_RIOT_API_KEY is not set")
	}

	if c.kevinPostgresConnection == "" {
		return errors.New("KEVIN_POSTGRES_CONNECTION is not set")
	}

	return nil
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) GetRiotAPIKey() string {
	return c.kevinRiotAPIKey
}

func (c *Config) GetPostgresConnection() string {
	return c.kevinPostgresConnection
}

// config is an unexported struct for unmarshaling from viper.
type config struct {
	KevinPostgresConnection string `mapstructure:"KEVIN_POSTGRES_CONNECTION"`
	KevinRiotAPIKey         string `mapstructure:"KEVIN_RIOT_API_KEY"`

	// env is the environment the application is running in.
	// Can be "dev" or "prod"; the default is "prod".
	Env string `mapstructure:"ENV"`

	// port is the port number to listen on.
	Port int `mapstructure:"PORT"`
}

func (c *config) ToConfig() *Config {
	return &Config{
		kevinPostgresConnection: c.KevinPostgresConnection,
		kevinRiotAPIKey:         c.KevinRiotAPIKey,
		env:                     c.Env,
		port:                    c.Port,
	}
}
