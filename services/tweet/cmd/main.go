package main

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/config"
	"github.com/Drivello/Twittah/services/tweet/ent"

	"github.com/Drivello/Twittah/services/tweet/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	common.InitLogger(cfg.LogLevel)
	logger := common.Logger()
	defer logger.Sync()

	// Inicializar Ent (Postgres)
	entClient, err := ent.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer entClient.Close()
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logger.Fatal("failed to run Ent migration", zap.Error(err))
	}

	// Inicializar Gin
	r := gin.Default()

	// Handlers principales

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	logger.Info("Tweet service started", zap.String("addr", cfg.Port))
	if err := r.Run(cfg.Port); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
