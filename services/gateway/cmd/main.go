package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gwhttp "github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	adapters "github.com/Drivello/Twittah/services/gateway/internal/adapters"
	"github.com/Drivello/Twittah/services/gateway/internal/common"
	"github.com/Drivello/Twittah/services/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Drivello/Twittah/services/gateway/config"
)

// main is the entry point for the Gateway Service.
// It loads configuration, sets up logging, initializes Kafka producers and HTTP handlers,
// and starts the Gin HTTP server with graceful shutdown support.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		common.Logger().Fatal("failed to load config", zap.Error(err))
	}
	common.InitLogger(cfg.LogLevel)
	defer common.Logger().Sync()

	authProducer, userProducer, tweetProducer, err := config.InitKafkaProducers(cfg)
	if err != nil {
		common.Logger().Fatal("failed to create kafka producers", zap.Error(err))
	}

	userServiceAdapter := adapters.NewUserServiceHTTPAdapter(cfg.UserMicroserviceURL)
	getFollowersUseCase := usecase.NewGetFollowersUseCase(userServiceAdapter)
	userQueryHandler := gwhttp.NewUserQueryHandler(getFollowersUseCase)

	registerUserUseCase := usecase.NewRegisterUserUseCase(authProducer)
	followUseCase := usecase.NewFollowUseCase(userProducer)
	unfollowUseCase := usecase.NewUnfollowUseCase(userProducer)
	createTweetUseCase := usecase.NewCreateTweetUseCase(tweetProducer)

	authHandler := gwhttp.NewAuthHandler(registerUserUseCase)
	userHandler := gwhttp.NewUserHandler(followUseCase, unfollowUseCase)
	tweetHandler := gwhttp.NewTweetHandler(createTweetUseCase)

	r := gin.Default()
	userQueryHandler.RegisterRoutes(r.Group("/api"))

	authGroup := r.Group("/auth")
	tweetGroup := r.Group("/tweets")
	userGroup := r.Group("/users")
	authHandler.RegisterRoutes(authGroup)
	tweetHandler.RegisterRoutes(tweetGroup)
	userHandler.RegisterRoutes(userGroup)

	r.GET("/health", gwhttp.NewHealthHandler(cfg).Health)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start HTTP server in goroutine
	go func() {
		common.Logger().Info("GatewayService corriendo", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Logger().Fatal("No se pudo iniciar el Gateway", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	common.Logger().Info("GatewayService: shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Logger().Error("GatewayService forced to shutdown", zap.Error(err))
	}
	common.Logger().Info("GatewayService exited cleanly")
}
