package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "photogallery/api_go/docs/swagger"
	"photogallery/api_go/internal/application/use_cases"
	"photogallery/api_go/internal/infrastructure/config"
	"photogallery/api_go/internal/infrastructure/persistence"
	"photogallery/api_go/internal/infrastructure/seed"
	infraServices "photogallery/api_go/internal/infrastructure/services"
	"photogallery/api_go/internal/infrastructure/uow"
	"photogallery/api_go/internal/web"
)

// @title Photo Gallery API
// @version 1.0
// @description API REST v1 Foto-tienda
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()
	db, err := persistence.NewDatabase(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	uowImpl := uow.NewUnitOfWork(db)
	jwtSvc := infraServices.NewJWTService(cfg.JWTSecret, cfg.JWTExpireSeconds)
	passSvc := infraServices.NewPasswordService()
	fileSvc := infraServices.NewLocalFileStorageService(cfg.FilesBasePath, cfg.FilesMaxSizeMB, cfg.FilesAllowedMIME)
	notifSvc := infraServices.NewNoopNotificationService()

	if err := seed.RunDevSeed(context.Background(), db, cfg, passSvc); err != nil {
		log.Fatalf("seed dev: %v", err)
	}

	uc := use_cases.NewUseCases(uowImpl, jwtSvc, passSvc, fileSvc, notifSvc)

	r := web.NewServer(uc, jwtSvc, cfg.CORSAllowedOrigins)
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
