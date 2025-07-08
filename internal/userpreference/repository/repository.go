package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-system/internal/userpreference/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"

)

type UserPreferenceRepository interface {
	Create(ctx context.Context, pref *model.UserPreference) error
	UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error
}

type userPreferenceRepo struct {
	coll *mongo.Collection
}

func NewUserPreferenceRepository(db *mongo.Database) UserPreferenceRepository {
	return &userPreferenceRepo{
		coll: db.Collection("UserPreferences"),
	}
}

func (r *userPreferenceRepo) Create(ctx context.Context, pref *model.UserPreference) error {
	_, err := r.coll.InsertOne(ctx, pref)
	return err
}

func (r *userPreferenceRepo) UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error {
	filter := bson.M{"user_id": updated.UserID}
	update := bson.M{"$set": updated}

	opts := options.Update().SetUpsert(false)
	result, err := r.coll.UpdateOne(ctx, filter, update,opts)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no preference found for user_id: %s", updated.UserID)
	}
	return nil
}
