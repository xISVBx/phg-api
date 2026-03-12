package main

import (
	"context"
	"log"

	"photogallery/api_go/internal/infrastructure/config"
	"photogallery/api_go/internal/infrastructure/persistence"
	"photogallery/api_go/internal/infrastructure/seed"
	infraServices "photogallery/api_go/internal/infrastructure/services"
)

func main() {
	cfg := config.Load()
	db, err := persistence.NewDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("seed-start: db connect: %v", err)
	}
	passSvc := infraServices.NewPasswordService()
	if err := seed.RunDevSeed(context.Background(), db, cfg, passSvc); err != nil {
		log.Fatalf("seed-start: run seed: %v", err)
	}
	log.Printf("seed-start: ejecutado")
}
