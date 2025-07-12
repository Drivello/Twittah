package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Drivello/Twittah/services/user/config"
	"github.com/Drivello/Twittah/services/user/ent"
	"github.com/Drivello/Twittah/services/user/internal/common"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthHandler struct {
	cfg       *config.Config
	Kafka     sarama.Client
	entClient *ent.Client
}

func NewHealthHandler(cfg *config.Config, db *ent.Client) *HealthHandler {
	kafkaClient, err := sarama.NewClient(cfg.KafkaBrokers, nil)
	if err != nil {
		zap.S().Fatal("failed to create kafka client for health handler", zap.Error(err))
	}
	return &HealthHandler{cfg: cfg, Kafka: kafkaClient, entClient: db}
}

func (h *HealthHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
}

func (h *HealthHandler) Health(c *gin.Context) {
	status := map[string]string{}
	ok := true

	requiredTopics := []string{h.cfg.KafkaUserConsumerConfig.Topic, h.cfg.KafkaUserConsumerConfig.DLQTopic}
	if h.Kafka != nil {
		topics, err := h.Kafka.Topics()
		topicsMap := map[string]bool{}
		for _, t := range topics {
			topicsMap[t] = true
		}
		if err != nil {
			status["kafka"] = err.Error()
			ok = false
		} else {
			missing := []string{}
			for _, t := range requiredTopics {
				if !topicsMap[t] {
					missing = append(missing, t)
				}
			}
			if len(missing) > 0 {
				status["kafka"] = "missing topics: " + strings.Join(missing, ", ")
				ok = false
			} else {
				status["kafka"] = "ok"
			}
		}
	}

	if h.entClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if _, err := h.entClient.User.Query().Limit(1).All(ctx); err != nil {
			status["database"] = err.Error()
			ok = false
		} else {
			status["database"] = "ok"
		}
		common.Logger().Debug("Database health check", zap.String("status", status["database"]))
	}

	if ok {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusInternalServerError, status)
	}
}
