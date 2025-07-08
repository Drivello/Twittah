package main

import (
	"context"

	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/internal/adapters/http"
	"github.com/Drivello/Twittah/services/user/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/user/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/user/internal/usecase"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := config.InitLogger()
	defer logger.Sync()

	cfg := config.LoadConfig()

	client, err := config.InitEntClient(cfg.PostgresDSN)
	if err != nil {
		logger.Sugar().Fatalw("failed opening connection to postgres", "error", err)
	}
	defer client.Close()

	brokers := config.ParseKafkaBrokers(cfg.KafkaBrokers)

	repo := postgres.NewEntUserRepository(client)
	followUC := usecase.NewFollowUserUseCase(repo, nil)
	unfollowUC := usecase.NewUnfollowUserUseCase(repo, nil)

	consumer, err := kafka.NewUserEventConsumer(brokers, "user-service-group", "follows", followUC, unfollowUC)
	if err != nil {
		logger.Sugar().Fatalw("failed to create kafka consumer", "error", err)
	}
	ctx := context.Background()
	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.Sugar().Fatalw("kafka consumer stopped", "error", err)
		}
	}()

	handler := http.NewUserHandler(followUC, unfollowUC)
	r := gin.Default()
	handler.RegisterRoutes(r)

	port := cfg.ServicePort
	if port == "" {
		port = "8082"
	}
	if err := r.Run(":" + port); err != nil {
		logger.Sugar().Fatalw("failed to run http server", "error", err)
	}
}

