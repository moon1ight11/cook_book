package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("environment", "development")
}
