package config

import (
	"time"

	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func ViperInit() {
	viper.AutomaticEnv()
}

// getEnvAsDuration parses a config key as a time.Duration or returns the default.
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := viper.GetString(key)
	if valStr != "" {
		val, err := time.ParseDuration(valStr)
		if err != nil {
			common.Logger().Error(
				"Invalid duration config value, using default",
				zap.String("key", key),
				zap.String("value", valStr),
				zap.Duration("default", defaultVal),
			)
			return defaultVal
		}
		return val
	}
	return defaultVal
}

// mustGetEnv returns the value or fatals if not set.
func mustGetEnv(key string) string {
	val := viper.GetString(key)
	if val == "" {
		common.Logger().Fatal(key + " env var/config required")
	}
	return val
}

// getEnvOrDefault returns the value or the default if not set.
func getEnvOrDefault(key, def string) string {
	val := viper.GetString(key)
	if val == "" {
		return def
	}
	return val
}
