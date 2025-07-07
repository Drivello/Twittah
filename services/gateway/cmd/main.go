package main

import (
	"github.com/Drivello/Twittah/services/gateway/config"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Configuración modular
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger, err := config.InitLogger()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	followProducer, userProducer, err := config.InitKafkaProducers(cfg)
	if err != nil {
		zap.S().Fatalw("failed to create kafka producers", "error", err)
	}

	r := gin.Default()

	authHandler := http.NewAuthHandler(userProducer)
	tweetHandler := http.NewTweetHandler(followProducer)
	followHandler := http.NewFollowHandler(followProducer)

	authGroup := r.Group("/auth")
	tweetGroup := r.Group("/tweets")
	followGroup := r.Group("") // root for follow endpoints

	authHandler.RegisterRoutes(authGroup)
	tweetHandler.RegisterRoutes(tweetGroup)
	followHandler.RegisterRoutes(followGroup)

	zap.S().Infow("GatewayService corriendo", "port", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		zap.S().Fatalw("No se pudo iniciar el Gateway", "error", err)
	}
}
