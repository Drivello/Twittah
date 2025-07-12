package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Drivello/Twittah/services/user/internal/common"
	"go.uber.org/zap"
)

// GetEnvAsDuration parses an environment variable as a time.Duration or returns the default.
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(key); valStr != "" {
		val, err := time.ParseDuration(valStr)
		if err != nil {
			common.Logger().Error("Invalid duration for %s: %s, using default %v", zap.String("key", key), zap.String("value", valStr), zap.Duration("default", defaultVal))
			return defaultVal
		}
		return val
	}
	return defaultVal
}

// mustGetEnv panics if the environment variable is not set.
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("env var %s required", key))
	}
	return val
}

// GetEnvOrDefault returns the environment variable value or the default if not set.
func getEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
