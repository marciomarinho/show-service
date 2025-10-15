package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stan-projects/show-service/internal/handlers"
)

func main() {
	router := gin.Default()
	router.POST("/shows", handlers.PostShows)

	// Define a simple GET endpoint
	router.GET("/health", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.Run()
}
