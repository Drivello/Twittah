package main

import (
	"os"

	"github.com/Drivello/Twittah/services/auth/config"
	"github.com/Drivello/Twittah/services/auth/ent"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/http"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/kafka"
	"github.com/Drivello/Twittah/services/auth/internal/adapters/postgres"
	"github.com/Drivello/Twittah/services/auth/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	zap.L().Info("[Startup] AuthService is starting up...")
	cfg := config.LoadConfig()

	producer, err := kafka.NewUserEventProducer(cfg.KafkaBrokers, "user_created")
	if err != nil {
		zap.L().Fatal("Failed to create Kafka producer", zap.Error(err))
	}

	// Instantiate repository and usecase layer
	entClient, err := ent.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}
	defer entClient.Close()

	repo := postgres.NewPostgresUserRepository(entClient)
	userUC := usecase.NewUserUseCasePort(repo)

	// Kafka consumers wiring is handled in kafka.StartKafkaConsumers for clarity
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
	handler := http.NewAuthHandler()
	handler.RegisterRoutes(r)

	zap.L().Info("AuthService listening on :%s", zap.String("port", cfg.Port))
	r.Run(":" + cfg.Port)
}
