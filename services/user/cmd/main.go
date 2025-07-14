package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	stdhttp "net/http"

	"github.com/Drivello/Twittah/services/user/config"
	userhttp "github.com/Drivello/Twittah/services/user/internal/adapters/http"
	"github.com/Drivello/Twittah/services/user/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/user/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/metrics"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	common.InitLogger(cfg.LogLevel)
	defer func() {
		if err := common.Logger().Sync(); err != nil {
			common.Logger().Error("Failed to sync logger", zap.Error(err))
		}
	}()

	entClient, err := config.InitEntClient(cfg.PostgresDSN)
	if err != nil {
		common.Logger().Fatalf("failed opening connection to postgres: %v", err)
	}
	defer entClient.Close()

	// Repositories
	repo := postgres.NewEntUserRepository(entClient)

	// Event Usecases
	createUserUC := usecase.NewCreateUserUseCase(repo)
	followUC := usecase.NewFollowUserUseCase(repo)
	unfollowUC := usecase.NewUnfollowUserUseCase(repo)

	// Read UseCases
	getFollowersUC := usecase.NewGetFollowersUseCase(repo)
	getFollowingUC := usecase.NewGetFollowingUseCase(repo)

	// Kafka
	producer, err := kafka.NewSyncProducer(cfg.KafkaBrokers)
	if err != nil {
		common.Logger().Fatalf("failed to start Kafka producer: %v", err)
	}
	defer producer.Close()

	// Worker pool
	workQueue := common.NewWorkQueue(true)
	kafkaEventDispatcher := kafka.NewKafkaEventDispatcher(createUserUC, followUC, unfollowUC)
	process := kafka.DefaultKafkaWorkProcess(kafkaEventDispatcher, producer, cfg.KafkaUserConsumerConfig.DLQTopic)
	go workQueue.Start(process, cfg.KafkaUserConsumerConfig.RetryConfig.MaxRetryDuration)
	defer workQueue.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go kafka.StartKafkaConsumers(ctx, cfg, repo, kafkaEventDispatcher, producer, workQueue)

	userHandler := userhttp.NewUserHandler(getFollowersUC, getFollowingUC)
	healthHandler := userhttp.NewHealthHandler(cfg, entClient)

	r := gin.Default()

	metrics.Init()
	r.GET("/metrics", gin.WrapH(metrics.Handler()))

	userHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	server := &stdhttp.Server{Addr: ":" + cfg.ServicePort, Handler: r}

	go func() {
		common.Logger().Info("[UserService] HTTP server starting", zap.String("port", cfg.ServicePort))
		if err := server.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			common.Logger().Fatalf("failed to run http server: %v", err)
		}
	}()

	common.Logger().Info("[UserService] Service started successfully and is ready to accept requests", zap.String("port", cfg.ServicePort))

	<-shutdown
	common.Logger().Info("[UserService] Shutdown signal received, shutting down gracefully...")
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()
	if err := server.Shutdown(ctxTimeout); err != nil {
		common.Logger().Error("[UserService] Error during server shutdown", zap.Error(err))
	}
	common.Logger().Info("[UserService] Server exited cleanly")
}
