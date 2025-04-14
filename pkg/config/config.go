package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Ollama  OllamaConfig  `mapstructure:"ollama"`
	Logging LoggingConfig `mapstructure:"logging"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type OllamaConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	DefaultModel string `mapstructure:"default_model"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		Ollama: OllamaConfig{
			BaseURL:      "http://localhost:11434",
			DefaultModel: "llama2",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
