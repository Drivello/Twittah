package http

import (
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Kafka sarama.Client
}

func (h *HealthHandler) Health(c *gin.Context) {
	status := map[string]string{}
	ok := true

	requiredTopics := []string{"tweets.published", "follows.created", "follows.deleted", "timelines.updated"}
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
