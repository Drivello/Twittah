package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/ent"
	authhttp "github.com/Drivello/Twittah/services/auth/internal/adapters/http"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/auth/internal/common"
	"github.com/Drivello/Twittah/services/auth/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	// Init logger
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	common.InitLogger(logLevel)
	logger := common.Logger()
	defer logger.Sync()

	// Load config
	cfg := config.LoadConfig()

	// Kafka producer
	producer, err := kafka.NewUserEventProducer(cfg.KafkaBrokers, cfg.KafkaUserEventsTopic)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}

	// Database connection
	entClient, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	repo := postgres.NewPostgresUserRepository(entClient)
	userUC := usecase.NewAuthUseCasesPort(repo, producer, cfg.UserCreatedEventType)

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Inicializa el worker pool y pásalo al consumer
	workerCount, queueSize := common.CalculateWorkerConfig(true)
	workQueue := common.NewWorkQueue(workerCount, queueSize)

	// Start Kafka consumers with cancelable context
	go kafka.StartKafkaConsumers(ctx, cfg, repo, userUC, producer.Producer, workQueue)

	// Setup Gin HTTP server
	r := gin.Default()
	handler := authhttp.NewAuthHandler(userUC)
	healthHandler := authhttp.NewHealthHandler(cfg, entClient)
	handler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

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

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-quit
	logger.Info("Shutdown signal received", zap.String("signal", sig.String()))

	// Cancel context for consumers and other goroutines
	cancel()

	// Shutdown HTTP server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	// Close DB connection
	if err := entClient.Close(); err != nil {
		logger.Error("Error closing database connection", zap.Error(err))
	}

	// Close Kafka producer
	kafka.CloseKafka(producer.Producer)

	logger.Info("AuthService exited cleanly")
}
