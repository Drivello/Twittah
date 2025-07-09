package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gwhttp "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	gwkafka "github.com/Drivello/Twittah/services/gateway/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/gateway/config"
	"github.com/IBM/sarama"
)

// main is the entry point for the Gateway Service.
// It loads configuration, sets up logging, initializes Kafka producers and HTTP handlers,
// and starts the Gin HTTP server with graceful shutdown support.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		common.Logger().Fatal("failed to load config", zap.Error(err))
	}
	logger := common.Logger()
	defer logger.Sync()

	followProducer, userProducer, err := config.InitKafkaProducers(cfg)
	if err != nil {
		logger.Fatal("failed to create kafka producers", zap.Error(err))
	}

	// Inicializa el SyncProducer de Kafka para tweets
	tweetSyncProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, nil)
	if err != nil {
		logger.Fatal("failed to create tweet kafka producer", zap.Error(err))
	}
	tweetProducer := gwkafka.NewTweetProducer(tweetSyncProducer, "tweets.published")
	tweetUsecase := usecase.NewTweetUsecase(tweetProducer)

	r := gin.Default()
	userUseCase := usecase.NewUserUseCase(userProducer)
	authHandler := gwhttp.NewAuthHandler(userUseCase)
	tweetHandler := gwhttp.NewTweetHandler(tweetUsecase)
	followUseCase := usecase.NewFollowUseCase(followProducer)
	followHandler := gwhttp.NewFollowHandler(followUseCase)

	authGroup := r.Group("/auth")
	tweetGroup := r.Group("/tweets")
	followGroup := r.Group("") // root for follow endpoints
	authHandler.RegisterRoutes(authGroup)
	tweetHandler.RegisterRoutes(tweetGroup)
	followHandler.RegisterRoutes(followGroup)
authHandler.RegisterRoutes(authGroup)
tweetHandler.RegisterRoutes(tweetGroup)
followHandler.RegisterRoutes(followGroup)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	// Start HTTP server in goroutine
	go func() {
		logger.Info("GatewayService corriendo", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("No se pudo iniciar el Gateway", zap.Error(err))
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
		logger.Error("GatewayService forced to shutdown", zap.Error(err))
	}
	logger.Info("GatewayService exited cleanly")
}

