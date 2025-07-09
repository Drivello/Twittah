package main

import (
	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/internal/adapters/http"
	"github.com/Drivello/Twittah/services/user/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/user/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger := common.Logger()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	client, err := config.InitEntClient(cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("failed opening connection to postgres", zap.Error(err))
	}
	defer client.Close()

	repo := postgres.NewEntUserRepository(client)
	followUC := usecase.NewFollowUserUseCase(repo, nil)
	unfollowUC := usecase.NewUnfollowUserUseCase(repo, nil)

	producer, err := kafka.StartAllKafkaConsumers(repo, cfg.KafkaBrokers, logger)
	if err != nil {
		logger.Fatal("failed to start Kafka consumers", zap.Error(err))
	}
	defer producer.Close()

	handler := http.NewUserHandler(followUC, unfollowUC)
	r := gin.Default()
	handler.RegisterRoutes(r)

	if err := r.Run(":" + cfg.ServicePort); err != nil {
		logger.Fatal("failed to run http server", zap.Error(err))
	}
}
