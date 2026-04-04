package config

import "time"

// общая конфигурация
type Config struct {
	Environment string       `mapstructure:"environment"`
	Server      ServerConfig `mapstructure:"server"`
	Mongo       MongoConfig  `mapstructure:"database"`
	Admin       AdminConfig  `mapstructure:"admin"`
	JWT         JWTConfig    `mapstructure:"jwt"`
	Logger      LoggerConfig `mapstructure:"logger"`
	Cache       CacheConfig  `mapstructure:"cache"`
}

// конфигурация сервера
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// конфигурация бд
type MongoConfig struct {
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}

// конфигурация админа сайта
type AdminConfig struct {
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Email    string `mapstructure:"email"`
}

// jwt-конфигурация
type JWTConfig struct {
	AccessSecret  string        `mapstructure:"access_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

// конфигурация логгера
type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"filepath"`
	MaxSize  int    `mapstructure:"maxsize"`
}

// конфигурация кэша
type CacheConfig struct {
	TTL             time.Duration `mapstructure:"ttl"`
	CleanUpInterval time.Duration `mapstructure:"cleanupinterval"`
}
