// Package config loads and validates TweetService configuration.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for TweetService.
type Config struct {
	PostgresDSN      string
	RedisAddr        string
	KafkaBrokers     []string
	KafkaGroupID     string
	ServicePort      string
	TimelineTTLHours int
}

// LoadConfig loads configuration from environment variables using Viper.
func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVICE_PORT", "8080")
	viper.SetDefault("TIMELINE_TTL_HOURS", 24)

	required := []string{"TWEET_POSTGRES_DSN", "REDIS_ADDR", "KAFKA_BROKERS", "KAFKA_GROUP_ID"}
	for _, key := range required {
		if viper.GetString(key) == "" {
			return nil, fmt.Errorf("missing required env var: %s", key)
		}
	}

	cfg := &Config{
		PostgresDSN:      viper.GetString("TWEET_POSTGRES_DSN"),
		RedisAddr:        viper.GetString("REDIS_ADDR"),
		KafkaBrokers:     viper.GetStringSlice("KAFKA_BROKERS"),
		KafkaGroupID:     viper.GetString("KAFKA_GROUP_ID"),
		ServicePort:      viper.GetString("SERVICE_PORT"),
		TimelineTTLHours: viper.GetInt("TIMELINE_TTL_HOURS"),
	}

	return cfg, nil
}
