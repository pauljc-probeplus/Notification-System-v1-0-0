package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-system/internal/scheduler/model"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"

	// "fmt"
	// "errors"

)

type SchedulerEntryRepository interface {
	CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error
}

type schedulerEntryRepo struct {
	coll *mongo.Collection
}

func NewSchedulerEntryRepository(db *mongo.Database) SchedulerEntryRepository {
	return &schedulerEntryRepo{
		coll: db.Collection("SchedulerEntry"),
	}
}

func (r *schedulerEntryRepo) CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error {
	_, err := r.coll.InsertOne(ctx, entry)
	return err
}

