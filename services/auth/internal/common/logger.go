package common

import (
	"go.uber.org/zap"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// Logger returns a singleton zap.Logger instance
func Logger() *zap.Logger {
	once.Do(func() {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = l
	})
	return logger
}
