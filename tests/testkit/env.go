package testkit

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"photogallery/api_go/internal/application/use_cases"
	"photogallery/api_go/internal/domain/entities"
	"photogallery/api_go/internal/infrastructure/persistence"
	infraServices "photogallery/api_go/internal/infrastructure/services"
	"photogallery/api_go/internal/infrastructure/uow"
	"photogallery/api_go/internal/web"
)

type IntegrationEnv struct {
	DB        *gorm.DB
	UoW       *uow.UnitOfWork
	UseCases  *use_cases.UseCases
	JWT       *infraServices.JWTService
	Router    *gin.Engine
	ActorUser *entities.User
}

func NewIntegrationEnv(t *testing.T) *IntegrationEnv {
	t.Helper()
	dsn := os.Getenv("DATABASE_DSN_TEST")
	if dsn == "" {
		t.Fatalf("DATABASE_DSN_TEST no definido (requerido para tests de integración)")
	}

	db, err := persistence.NewDatabase(dsn)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("open sql db: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping test db: %v", err)
	}
	db = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
	ResetDatabase(t, db)

	actor := &entities.User{
		Username:     "test_actor_" + uuid.NewString()[:8],
		PasswordHash: "test_hash",
		FullName:     "Test Actor",
		Email:        "test.actor+" + time.Now().UTC().Format("150405.000000") + "@example.com",
		IsActive:     true,
	}
	if err := db.Create(actor).Error; err != nil {
		t.Fatalf("seed actor user: %v", err)
	}

	jwtSvc := infraServices.NewJWTService("test-secret", 3600)
	passSvc := infraServices.NewPasswordService()
	fsSvc := infraServices.NewLocalFileStorageService(t.TempDir(), 25, []string{"image/jpeg", "image/png", "application/pdf"})
	notifSvc := infraServices.NewNoopNotificationService()

	uowImpl := uow.NewUnitOfWork(db)
	uc := use_cases.NewUseCases(uowImpl, jwtSvc, passSvc, fsSvc, notifSvc)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	r := web.NewServer(uc, jwtSvc, []string{"http://localhost:5173"})

	return &IntegrationEnv{DB: db, UoW: uowImpl, UseCases: uc, JWT: jwtSvc, Router: r, ActorUser: actor}
}

func (e *IntegrationEnv) AuthHeader(t *testing.T, userID uuid.UUID, username string) string {
	t.Helper()
	tokens, err := e.JWT.Generate(userID, username)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return "Bearer " + tokens.AccessToken
}

func ResetDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()
	sql := `
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'TRUNCATE TABLE public.' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
    END LOOP;
END
$$;`
	if err := db.Exec(sql).Error; err != nil {
		t.Fatalf("reset test db: %v", err)
	}
}
