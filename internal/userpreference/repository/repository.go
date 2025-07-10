package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"notification-system/internal/userpreference/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"
	"errors"

)

type UserPreferenceRepository interface {
	Create(ctx context.Context, pref *model.UserPreference) error
	UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error
	CountByFilter(ctx context.Context, filter bson.M) (int64, error) //get all existing user ids
	GetUserPreference(ctx context.Context, userID string) (*model.UserPreference, error)

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
	if err != nil {
        // Check for duplicate key error
        if mongo.IsDuplicateKeyError(err) {
            return errors.New("user preference already exists")
        }
        return err
    }
    return nil
	
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

func (r *userPreferenceRepo) CountByFilter(ctx context.Context, filter bson.M) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}

func (r *userPreferenceRepo) GetUserPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	filter := bson.M{"user_id": userID}
	var pref model.UserPreference
	err := r.coll.FindOne(ctx, filter).Decode(&pref)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no preference found for user_id: %s", userID)
		}
		return nil, err
	}
	return &pref, nil
}
