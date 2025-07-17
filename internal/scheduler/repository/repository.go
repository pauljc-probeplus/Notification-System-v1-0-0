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
	LogFailure(ctx context.Context, log *model.FailureLog) error
}

type schedulerEntryRepo struct {
	coll *mongo.Collection
	failureColl *mongo.Collection
}

func NewSchedulerEntryRepository(db *mongo.Database) SchedulerEntryRepository {
	return &schedulerEntryRepo{
		coll: db.Collection("SchedulerEntry"),
		failureColl: db.Collection("FailureLogs"),
	}
}

func (r *schedulerEntryRepo) CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error {
	_, err := r.coll.InsertOne(ctx, entry)
	return err
}

func (r *schedulerEntryRepo) LogFailure(ctx context.Context, log *model.FailureLog) error {
	_, err := r.failureColl.InsertOne(ctx, log)
	return err
}
