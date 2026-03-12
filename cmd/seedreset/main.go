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
		log.Fatalf("seed-reset: db connect: %v", err)
	}

	ctx := context.Background()
	if err := seed.ResetDevSeedData(ctx, db, cfg); err != nil {
		log.Fatalf("seed-reset: reset dev data: %v", err)
	}

	if !cfg.SeedEnabled {
		log.Printf("seed-reset: reset completado. SEED_ENABLED=false, se omite recarga.")
		return
	}

	passSvc := infraServices.NewPasswordService()
	if err := seed.RunDevSeed(ctx, db, cfg, passSvc); err != nil {
		log.Fatalf("seed-reset: run seed: %v", err)
	}
	log.Printf("seed-reset: reset y recarga completados")
}
