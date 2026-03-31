package config

import "time"

type Config struct {
	Environment string       `mapstructure:"environment"`
	Server      ServerConfig `mapstructure:"server"`
	Mongo       MongoConfig  `mapstructure:"mongo"`
	JWT         JWTConfig    `mapstructure:"jwt"`
	Logger      LoggerConfig `mapstructure:"logger"`
	Cache       CacheConfig  `mapstructure:"cache"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type MongoConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filepath"`
	MaxSize  int    `mapstructure:"maxsize"`
}

type CacheConfig struct {
	TTL             time.Duration `mapstructure:"ttl"`
	CleanUpInterval time.Duration `mapstructure:"cleanupinterval"`
}
