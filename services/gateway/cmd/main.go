package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gwhttp "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/gateway/config"
)

// main is the entry point for the Gateway Service.
// It loads configuration, sets up logging, initializes Kafka producers and HTTP handlers,
// and starts the Gin HTTP server with graceful shutdown support.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.S().Fatalw("failed to load config", "error", err)
	}

	logger, err := config.InitLogger()
	if err != nil {
		zap.S().Fatalw("failed to init logger", "error", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	followProducer, userProducer, err := config.InitKafkaProducers(cfg)
	if err != nil {
		zap.S().Fatalw("failed to create kafka producers", "error", err)
	}

	r := gin.Default()
	userUseCase := usecase.NewUserUseCase(userProducer)
authHandler := gwhttp.NewAuthHandler(userUseCase)
	tweetHandler := gwhttp.NewTweetHandler(followProducer) // TODO: Refactor to use usecase if tweet logic is added
	followUseCase := usecase.NewFollowUseCase(followProducer)
followHandler := gwhttp.NewFollowHandler(followUseCase)

authGroup := r.Group("/auth")
tweetGroup := r.Group("/tweets")
followGroup := r.Group("") // root for follow endpoints
authHandler.RegisterRoutes(authGroup)
tweetHandler.RegisterRoutes(tweetGroup)
followHandler.RegisterRoutes(followGroup)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	// Start HTTP server in goroutine
	go func() {
		zap.S().Infow("GatewayService corriendo", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalw("No se pudo iniciar el Gateway", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	zap.S().Info("GatewayService: shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Errorw("GatewayService forced to shutdown", "error", err)
	}
	zap.S().Info("GatewayService exited cleanly")
}

