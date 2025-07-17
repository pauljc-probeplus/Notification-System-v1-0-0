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
	LogFailure(ctx context.Context, log *model.FailureLog) error 
}

type schedulerEntryService struct {
	repo repository.SchedulerEntryRepository
	//logRepo repository.FailureLogRepository
}

func NewSchedulerEntryService(repo repository.SchedulerEntryRepository) SchedulerEntryService {
	return &schedulerEntryService{
		repo: repo,
		//logRepo: logRepo,
	}
}

func (s *schedulerEntryService) CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error {
	
	return s.repo.CreateSchedulerEntry(ctx, entry)
}

func (s *schedulerEntryService) LogFailure(ctx context.Context, log *model.FailureLog) error {
	return s.repo.LogFailure(ctx, log)
}



