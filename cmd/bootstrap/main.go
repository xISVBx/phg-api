package main

import (
	"context"
	"log"

	"photogallery/api_go/internal/infrastructure/bootstrap"
	"photogallery/api_go/internal/infrastructure/config"
	"photogallery/api_go/internal/infrastructure/persistence"
	infraServices "photogallery/api_go/internal/infrastructure/services"
)

func main() {
	cfg := config.Load()
	db, err := persistence.NewDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("bootstrap: db connect: %v", err)
	}

	svc := bootstrap.NewService(db, cfg, infraServices.NewPasswordService())
	result, err := svc.RunManual(context.Background())
	if err != nil {
		log.Fatalf("bootstrap: run manual: %v", err)
	}

	log.Printf("bootstrap: finalizado trigger=%s skipped=%t reason=%s", result.Trigger, result.Skipped, result.Reason)
}
