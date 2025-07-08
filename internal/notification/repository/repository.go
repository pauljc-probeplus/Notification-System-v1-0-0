package repository

import (
	"context"
	"notification-system/internal/notification/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
}

type notificationRepo struct {
	coll *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepo{
		coll: db.Collection("Notifications"),
	}
}

func (r *notificationRepo) Create(ctx context.Context, n *model.Notification) error {
	_, err := r.coll.InsertOne(ctx, n)
	return err
}
