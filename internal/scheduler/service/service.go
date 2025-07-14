package service

import (
	"context"
	"notification-system/internal/scheduler/model"
	"notification-system/internal/scheduler/repository"
	//"go.mongodb.org/mongo-driver/bson"
	//"fmt"
)

type SchedulerEntryService interface {
	CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error
}

type schedulerEntryService struct {
	repo repository.SchedulerEntryRepository
}

func NewSchedulerEntryService(repo repository.SchedulerEntryRepository) SchedulerEntryService {
	return &schedulerEntryService{repo: repo}
}

func (s *schedulerEntryService) CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error {
	
	return s.repo.CreateSchedulerEntry(ctx, entry)
}




