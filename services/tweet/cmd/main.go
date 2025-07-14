package main

import (
	"context"

	"github.com/Drivello/Twittah/services/tweet/config"

	"github.com/Drivello/Twittah/services/tweet/internal/adapters"
	"github.com/Drivello/Twittah/services/tweet/internal/adapters/http"
	"github.com/Drivello/Twittah/services/tweet/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/tweet/internal/adapters/postgres"
	redisadapter "github.com/Drivello/Twittah/services/tweet/internal/adapters/redis"
	"github.com/Drivello/Twittah/services/tweet/internal/common"
	"github.com/Drivello/Twittah/services/tweet/internal/metrics"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	common.InitLogger(cfg.LogLevel)
	defer func() {
		if err := common.Logger().Sync(); err != nil {
			common.Logger().Error("Failed to sync logger", zap.Error(err))
		}
	}()

	// Inicializar Ent (Postgres)
	entClient, err := config.InitEntClient(cfg.PostgresDSN)
	if err != nil {
		common.Logger().Fatal("failed to connect to database", zap.Error(err))
	}
	defer entClient.Close()

	// Services
	userService := adapters.NewUserServiceHTTPAdapter(cfg.UserServiceURL)

	// Repositories
	tweetRepo := postgres.NewTweetRepository(entClient)
	userRepo := postgres.NewUserRepository(entClient)

	// Event Usecases
	createUserUC := usecase.NewCreateUserUsecase(userRepo)
	createTweetUC := usecase.NewCreateTweetUsecase(tweetRepo)
	deleteTweetUC := usecase.NewDeleteTweetUsecase(tweetRepo)
	// Redis Timeline Cache
	redisClient := redisadapter.NewRedisClient(cfg.RedisAddr)
	timelineCache := redisadapter.NewTimelineCache(redisClient, cfg.RedisTimelineTTLHours)

	getTimelineUC := usecase.NewGetTimelineUseCase(tweetRepo, userService, timelineCache)
	getTweetsFromMultipleUserIDsUC := usecase.NewGetTweetsFromMultipleUserIDsUsecase(tweetRepo)
	getUserTweetsUC := usecase.NewGetUserTweetsUsecase(tweetRepo)

	tweetQueryHandler := http.NewTweetQueryHandler(getTimelineUC, getTweetsFromMultipleUserIDsUC, getUserTweetsUC)

	// Inicializar Kafka
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

	kafkaEventDispatcher := kafka.NewKafkaEventDispatcher(createUserUC, createTweetUC, deleteTweetUC)

	kafka.StartKafkaConsumers(ctx, cfg, kafkaEventDispatcher, producer, workQueue)

	// Inicializar Gin
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(metrics.Handler()))

	tweetGroup := r.Group("/")
	tweetQueryHandler.RegisterRoutes(tweetGroup)

	common.Logger().Info("Tweet service started", zap.String("addr", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		common.Logger().Fatal("server failed", zap.Error(err))
	}
}
