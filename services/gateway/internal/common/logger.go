package common

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger   *zap.Logger
	once     sync.Once
	logLevel string = "info"
)

// InitLogger sets the log level and resets the logger singleton
func InitLogger(level string) {
	logLevel = level
	logger = nil
	once = sync.Once{} // reset singleton
}

// Logger returns a singleton zap.Logger instance with the configured level
func Logger() *zap.SugaredLogger {
	once.Do(func() {
		cfg := zap.NewProductionConfig()
		if err := cfg.Level.UnmarshalText([]byte(logLevel)); err != nil {
			cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		}
		l, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		logger = l
	})
	return logger.Sugar()
}
