package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("COOKBOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

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

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")

	viper.SetDefault("admin.name", "")
	viper.SetDefault("admin.password", "")
	viper.SetDefault("admin.email", "")

	viper.SetDefault("database.uri", "")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 27017)
	viper.SetDefault("database.name", "cook_book")
	viper.SetDefault("database.user", "")
	viper.SetDefault("database.password", "")

	viper.SetDefault("jwt.access_secret", "")
	viper.SetDefault("jwt.refresh_secret", "")
	viper.SetDefault("jwt.access_expiry", 15*time.Minute)
	viper.SetDefault("jwt.refresh_expiry", 24*time.Hour)

	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.filepath", "logs/app.log")
	viper.SetDefault("logger.maxsize", 10)

	viper.SetDefault("cache.ttl", 5*time.Minute)
	viper.SetDefault("cache.cleanupinterval", time.Minute)
}
