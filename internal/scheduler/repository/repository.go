package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-system/internal/scheduler/model"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"

	// "fmt"
	// "errors"
	"notification-system/internal/common/timeutil"

)

type SchedulerEntryRepository interface {
	CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error
	LogFailure(ctx context.Context, log *model.FailureLog) error
	FetchPendingEntries(ctx context.Context) ([]model.SchedulerEntry, error)
	MarkAsProcessed(ctx context.Context, id string) error
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

//for kafka
func (r *schedulerEntryRepo) FetchPendingEntries(ctx context.Context) ([]model.SchedulerEntry, error) {
	//collection := r.db.Collection("scheduler_entry")
	cursor, err := r.coll.Find(ctx, bson.M{
		"status": "pending",
		"send_at": bson.M{
			"$lte": timeutil.NowISTFormatted(),  // Use $lte operator here
		},
	})
	if err != nil {
		return nil, err
	}
	var entries []model.SchedulerEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

//for kafka
func (r *schedulerEntryRepo) MarkAsProcessed(ctx context.Context, id string) error {
	//collection := r.db.Collection("scheduler_entry")
	_, err := r.coll.UpdateOne(ctx, bson.M{"scheduler_entry_id": id}, bson.M{"$set": bson.M{"status": "processed"}})
	return err
}

