package config

import (
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/spf13/viper"
)

func ViperInit() {
	viper.AutomaticEnv()
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
