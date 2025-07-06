package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/Drivello/Twittah/services/gateway/internal/adapters/http"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	handler := http.NewGatewayHandler()
	// TODO: Inject real usecases/repos when ready

	handler.RegisterRoutes(r)

	log.Printf("GatewayService corriendo en :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("No se pudo iniciar el Gateway: %v", err)
	}
}
