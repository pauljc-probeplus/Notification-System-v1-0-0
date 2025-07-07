package service

import (
	"context"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification/repository"
	//"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, n *notification.Notification) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(ctx context.Context, n *notification.Notification) error {
	// TODO: Add business logic for DND, limits, etc later
	return s.repo.Create(ctx, n)
}
