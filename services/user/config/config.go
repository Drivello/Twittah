package config

import (
	"log"
	"github.com/spf13/viper"
)

// Config holds all configuration for the UserService.
type Config struct {
	PostgresDSN  string
	RedisAddr    string
	KafkaBrokers []string
	ServicePort  string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, relying on environment variables")
	}

	cfg := &Config{
		PostgresDSN:  viper.GetString("POSTGRES_DSN"),
		RedisAddr:    viper.GetString("REDIS_ADDR"),
		KafkaBrokers: viper.GetStringSlice("KAFKA_BROKERS"),
		ServicePort:  viper.GetString("USER_SERVICE_PORT"),
	}

	if cfg.ServicePort == "" {
		cfg.ServicePort = "8082"
	}
	return cfg
}
