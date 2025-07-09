// Command main is the entrypoint for TweetService.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Drivello/Twittah/services/tweet/config"
	// "github.com/Drivello/Twittah/services/tweet/ent" // Disabled: ent package not generated
	tweethttp "github.com/Drivello/Twittah/services/tweet/internal/adapters/http"
	"github.com/Drivello/Twittah/services/tweet/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/tweet/internal/adapters/postgres"
	redisadapter "github.com/Drivello/Twittah/services/tweet/internal/adapters/redis"
	redis "github.com/go-redis/redis/v8"
	"github.com/Drivello/Twittah/services/tweet/internal/usecase"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Sugar().Fatalw("failed to load config", "error", err)
		os.Exit(1)
	}

	// Wire Ent client
	// TODO: Enable Ent ORM when ent is generated
	// entClient, err := ent.Open("postgres", cfg.PostgresDSN)
	// if err != nil {
	// 	logger.Sugar().Fatalw("failed to connect to postgres", "error", err)
	// 	os.Exit(1)
	// }
	// defer // entClient.Close() // Disabled: ent not generated
	// repo := postgres.NewEntTweetRepository(entClient)
	repo := postgres.NewEntTweetRepository(nil)

	// Wire Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()
	cache := redisadapter.NewTimelineCache(redisClient, cfg.TimelineTTLHours)

	publishUC := usecase.NewPublishTweet(repo, cache)
	timelineUC := usecase.NewGetTimeline(repo, cache)

	// Wire Kafka event handlers (tweet published, follows created/deleted, timelines updated)
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaCfg)
	if err != nil {
		logger.Sugar().Fatalw("failed to create kafka producer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Sugar().Errorw("failed to close kafka producer", "error", err)
		}
	}()

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, sarama.NewConfig())
	if err != nil {
		logger.Sugar().Fatalw("failed to create kafka consumer group", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			logger.Sugar().Errorw("failed to close kafka consumer group", "error", err)
		}
	}()

	// Instantiate real Kafka adapters
	producerAdapter := kafka.NewTimelineProducer(producer, logger)
	consumerHandler := kafka.NewConsumerGroupHandler(publishUC, timelineUC, cache, logger)

	go func() {
		for {
			if err := consumerGroup.Consume(context.Background(), []string{"tweets.published", "follows.created", "follows.deleted"}, consumerHandler); err != nil {
				logger.Sugar().Errorw("kafka consume error", "error", err)
				break
			}
		}
	}()

	_ = producerAdapter // Avoid unused warning; wire into usecases as needed

	router := gin.Default()
	handler := tweethttp.NewTweetHandler(publishUC, timelineUC)
	handler.RegisterRoutes(router)

	log.Printf("TweetService starting on port %s", cfg.ServicePort)

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: router,
	}
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatalw("http server error", "error", err)
		}
	}()
	<-shutdownCtx.Done()
	logger.Info("shutting down TweetService...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Sugar().Errorw("graceful shutdown failed", "error", err)
	}

	if err := redisClient.Close(); err != nil {
		logger.Sugar().Errorw("failed to close redis client", "error", err)
	}
	// TODO: Add Kafka wiring and shutdown
}
