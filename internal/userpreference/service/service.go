package service

import (
	"context"
	"notification-system/internal/userpreference/model"
	"notification-system/internal/userpreference/repository"
	"go.mongodb.org/mongo-driver/bson"
	"fmt"
)

type UserPreferenceService interface {
	CreateUserPreference(ctx context.Context, pref *model.UserPreference) error
	UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error
	GetUserPreference(ctx context.Context, userID string) (*model.UserPreference, error)

}

type userPreferenceService struct {
	repo repository.UserPreferenceRepository
}

func NewUserPreferenceService(repo repository.UserPreferenceRepository) UserPreferenceService {
	return &userPreferenceService{repo: repo}
}

func (s *userPreferenceService) CreateUserPreference(ctx context.Context, pref *model.UserPreference) error {
	// Check if user_id already exists
	filter := bson.M{"user_id": pref.UserID}
	count, err := s.repo.CountByFilter(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking user_id uniqueness: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("exists") // use a constant if needed
	}
	return s.repo.Create(ctx, pref)
}


func (s *userPreferenceService) UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error {
	return s.repo.UpdateUserPreference(ctx, updated)
}

func (s *userPreferenceService) GetUserPreference(ctx context.Context, userID string) (*model.UserPreference, error) {
	return s.repo.GetUserPreference(ctx, userID)
}

