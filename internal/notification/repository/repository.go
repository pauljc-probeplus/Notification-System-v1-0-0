package repository

import (
	"context"
	"github.com/pauljc-probeplus/Notification-System-v1-0-0/internal/notification"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *notification.Notification) error
}

type notificationRepo struct {
	coll *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepo{
		coll: db.Collection("notifications"),
	}
}

func (r *notificationRepo) Create(ctx context.Context, n *notification.Notification) error {
	_, err := r.coll.InsertOne(ctx, n)
	return err
}
