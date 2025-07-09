package config

import (
	"os"

	"go.uber.org/zap"
)

// InitZapLogger initializes zap logger with the log level from LOG_LEVEL env var (default info)
func InitZapLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	var lvl zap.AtomicLevel
	switch logLevel {
	case "debug":
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	cfgZap := zap.Config{
		Level:            lvl,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfgZap.Build()
	if err != nil {
		panic("failed to initialize zap logger")
	}
	zap.ReplaceGlobals(logger)
	zap.S().Infof("[Startup] AuthService is starting up with log level: %s", logLevel)
}
