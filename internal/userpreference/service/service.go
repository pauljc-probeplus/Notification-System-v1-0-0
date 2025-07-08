package service

import (
	"context"
	"notification-system/internal/userpreference/model"
	"notification-system/internal/userpreference/repository"
)

type UserPreferenceService interface {
	CreateUserPreference(ctx context.Context, pref *model.UserPreference) error
	UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error
}

type userPreferenceService struct {
	repo repository.UserPreferenceRepository
}

func NewUserPreferenceService(repo repository.UserPreferenceRepository) UserPreferenceService {
	return &userPreferenceService{repo: repo}
}

func (s *userPreferenceService) CreateUserPreference(ctx context.Context, pref *model.UserPreference) error {
	return s.repo.Create(ctx, pref)
}

func (s *userPreferenceService) UpdateUserPreference(ctx context.Context, updated *model.UserPreference) error {
	return s.repo.UpdateUserPreference(ctx, updated)
}
