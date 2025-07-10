package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/ent"
	authhttp "github.com/Drivello/Twittah/services/auth/internal/adapters/http"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/auth/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"go.uber.org/zap"
)

// main is the entry point for the AuthService. It sets up logging, configuration, database, Kafka, and the HTTP server.
// Implements graceful shutdown for HTTP server and resources.
func main() {
	logger := common.Logger()
defer logger.Sync()

	cfg := config.LoadConfig()

	producer, err := kafka.NewUserEventProducer(cfg.KafkaBrokers, "user_created")
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}

	entClient, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	repo := postgres.NewPostgresUserRepository(entClient)
	userUC := usecase.NewUserUseCasePort(repo)

	userConsumer := &kafka.UserCreateConsumer{
		UserUC:   userUC,
		Producer: producer.Producer,
		Config:   cfg.Retry,
	}
	dlqWorker := &kafka.DLQWorker{
		Repo:     repo,
		Producer: producer.Producer,
		Config:   cfg.DLQWorker,
	}
	kafka.StartKafkaConsumers(cfg, userConsumer, dlqWorker)

	r := gin.Default()
	handler := authhttp.NewAuthHandler()
	handler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start HTTP server in goroutine
	go func() {
		logger.Info("AuthService listening", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down AuthService...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}
	if err := entClient.Close(); err != nil {
		logger.Error("Error closing database connection", zap.Error(err))
	}
	// TODO: Add Kafka producer/consumer shutdown if needed

	logger.Info("AuthService exited cleanly")
}
