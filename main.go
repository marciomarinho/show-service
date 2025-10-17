package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/marciomarinho/show-service/internal/config"
	"github.com/marciomarinho/show-service/internal/database"
	"github.com/marciomarinho/show-service/internal/handlers"
	"github.com/marciomarinho/show-service/internal/repository"
	"github.com/marciomarinho/show-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if cfg.Env == config.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}

	// Infra
	dyn, err := database.NewDynamo(context.Background(), cfg)
	if err != nil {
		log.Fatalf("dynamo: %v", err)
	}

	// Repo
	repo := repository.NewShowRepository(dyn)

	// App
	svc := service.NewShowService(repo)

	// HTTP
	h := handlers.NewShowHandler(svc)
	r := gin.Default()

	// Apply authentication middleware for non-local environments
	r.Use(handlers.AuthMiddleware(cfg))

	// Health check endpoint (no auth required)
	r.GET("/health", handlers.HealthCheck)

	// Protected endpoints
	r.POST("/shows", h.PostShows)
	r.GET("/shows", h.GetShows)

	port := 8080
	log.Printf("env=%s table=%s listening=:%d", cfg.Env, dyn.TableName(), port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
