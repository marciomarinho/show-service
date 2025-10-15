package main

import (
	"github.com/gin-gonic/gin"

	"github.com/marciomarinho/show-service/internal/handlers"
)

func main() {
	router := gin.Default()
	router.GET("/health", handlers.HealthCheck)
	router.POST("/shows", handlers.PostShows)
	router.GET("/shows", handlers.GetShows)

	router.Run()
}
