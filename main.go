package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/marciomarinho/show-service/internal/config"
	"github.com/marciomarinho/show-service/internal/handlers"
	"github.com/marciomarinho/show-service/internal/infra"
	"github.com/marciomarinho/show-service/internal/repository"
	"github.com/marciomarinho/show-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if cfg.Env == config.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// Infra
	dyn, err := infra.NewDynamo(context.Background(), cfg)
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
	r.GET("/health", handlers.HealthCheck)
	r.POST("/shows", h.PostShows)
	r.GET("/shows", h.GetShows)
	r.GET("/shows/:slug", h.GetShowBySlug)

	port := 8080
	log.Printf("env=%s table=%s listening=:%d", cfg.Env, dyn.TableName, port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
