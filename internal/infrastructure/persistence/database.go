package persistence

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"photogallery/api_go/internal/domain/entities"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	if err := autoMigrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;").Error; err != nil {
		return fmt.Errorf("create extension pgcrypto: %w", err)
	}
	return db.AutoMigrate(
		&entities.User{}, &entities.Role{}, &entities.Menu{}, &entities.SubMenu{}, &entities.Permission{},
		&entities.RoleSubMenuPermission{}, &entities.UserRole{}, &entities.UserPermissionOverride{}, &entities.AuditLog{},
		&entities.Category{}, &entities.Product{}, &entities.Customer{}, &entities.Sale{}, &entities.SaleItem{}, &entities.SalePayment{},
		&entities.WorkOrder{}, &entities.WorkOrderItem{}, &entities.Appointment{}, &entities.File{}, &entities.FileLink{},
		&entities.Adjustment{}, &entities.AdjustmentItem{}, &entities.CashCategory{}, &entities.CashSession{}, &entities.CashMovement{},
		&entities.Worker{}, &entities.CommissionEntry{}, &entities.WorkerPayment{}, &entities.WorkerPaymentAllocation{},
		&entities.AppSetting{},
	)
}
