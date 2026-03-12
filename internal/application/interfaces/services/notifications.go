package services

import "context"

type NotificationPayload struct {
	Channel string
	To      string
	Title   string
	Body    string
}

type INotificationService interface {
	Send(ctx context.Context, payload NotificationPayload) error
}
