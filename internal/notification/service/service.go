package service

import (
	"context"
	"notification-system/internal/notification/model"
	"notification-system/internal/notification/repository"
	
)

type NotificationService interface {
	CreateNotification(ctx context.Context, n *model.Notification) error
}

type notificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(ctx context.Context, n *model.Notification) error {
	// TODO: Add business logic for DND, limits, etc later
	return s.repo.Create(ctx, n)
}
