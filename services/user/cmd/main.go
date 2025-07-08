package main

import (
	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/internal/adapters/http"
	"github.com/Drivello/Twittah/services/user/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/user/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/user/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	producer, err := kafka.StartAllKafkaConsumers(repo, brokers, logger)
	if err != nil {
		logger.Sugar().Fatalw("failed to start Kafka consumers", "error", err)
	}
	defer producer.Close()

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
