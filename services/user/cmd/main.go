package main

import (
	"context"

	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/internal/adapters/http"
	"github.com/Drivello/Twittah/services/user/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/user/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/Drivello/Twittah/services/user/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	common.InitLogger(cfg.LogLevel)
	defer common.Logger().Sync()

	client, err := config.InitEntClient(cfg.PostgresDSN)
	if err != nil {
		common.Logger().Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Migración automática Ent
	if err := client.Schema.Create(context.Background()); err != nil {
		common.Logger().Fatalf("failed to run Ent migration: %v", err)
	}

	// Repositories
	repo := postgres.NewEntUserRepository(client)

	// Event Usecases
	createUserUC := usecase.NewCreateUserUseCase(repo)
	followUC := usecase.NewFollowUserUseCase(repo, nil)
	unfollowUC := usecase.NewUnfollowUserUseCase(repo, nil)

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
	defer workQueue.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaEventDispatcher := kafka.NewKafkaEventDispatcher(createUserUC, followUC, unfollowUC)

	kafka.StartKafkaConsumers(ctx, cfg, repo, kafkaEventDispatcher, producer, workQueue)

	handler := http.NewUserHandler(getFollowersUC, getFollowingUC)
	r := gin.Default()
	handler.RegisterRoutes(r)

	healthHandler := http.NewHealthHandler(cfg, client)
	healthHandler.RegisterRoutes(r)

	if err := r.Run(":" + cfg.ServicePort); err != nil {
		common.Logger().Fatalf("failed to run http server: %v", err)
	}
}
