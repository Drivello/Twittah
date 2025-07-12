package http

import (
	"net/http"
	"strings"

	"github.com/Drivello/Twittah/services/gateway/config"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthHandler struct {
	cfg   *config.GatewayConfig
	Kafka sarama.Client
}

func NewHealthHandler(cfg *config.GatewayConfig) *HealthHandler {
	// Inicializa el health handler con el cliente de Kafka
	kafkaClient, err := sarama.NewClient(cfg.KafkaBrokers, nil)
	if err != nil {
		zap.S().Fatal("failed to create kafka client for health handler", zap.Error(err))
	}
	return &HealthHandler{cfg: cfg, Kafka: kafkaClient}
}

func (h *HealthHandler) Health(c *gin.Context) {
	status := map[string]string{}
	ok := true

	requiredTopics := []string{h.cfg.KafkaAuthTopic, h.cfg.KafkaUserTopic, h.cfg.KafkaTweetTopic}
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

	if ok {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusInternalServerError, status)
	}
}
