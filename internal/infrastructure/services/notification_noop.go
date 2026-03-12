package services

import (
	"context"

	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type NoopNotificationService struct{}

func NewNoopNotificationService() *NoopNotificationService { return &NoopNotificationService{} }

func (s *NoopNotificationService) Send(_ context.Context, _ appsvc.NotificationPayload) error {
	// TODO: conectar proveedor real (WhatsApp/email) usando este puerto.
	return nil
}
