package config

import (
	"log"
	"context"
	"strings"

	"go.uber.org/zap"
	"github.com/Drivello/Twittah/services/user/ent"
	"github.com/spf13/viper"
)

// Config holds all configuration for the UserService.
type Config struct {
	PostgresDSN  string
	RedisAddr    string
	KafkaBrokers []string
	ServicePort  string
}

func InitLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	return logger
}

func InitEntClient(postgresDSN string) (*ent.Client, error) {
	client, err := ent.Open("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func ParseKafkaBrokers(brokers []string) []string {
	if len(brokers) == 1 && strings.Contains(brokers[0], ",") {
		return strings.Split(brokers[0], ",")
	}
	return brokers
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
