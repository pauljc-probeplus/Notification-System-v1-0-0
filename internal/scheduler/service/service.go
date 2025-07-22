package service

import (
	"context"
	"notification-system/internal/scheduler/model"
	"notification-system/internal/scheduler/repository"
	"time"
	"log"
	"encoding/json"
	//"go.mongodb.org/mongo-driver/bson"
	
	//"fmt"
)

type SchedulerEntryService interface {
	CreateSchedulerEntry(ctx context.Context, entry *model.SchedulerEntry) error
	LogFailure(ctx context.Context, log *model.FailureLog) error 
	PollAndDispatch(ctx context.Context) error
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

//kafka logic
func (s *schedulerEntryService) PollAndDispatch(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			entries, err := s.repo.FetchPendingEntries(ctx)
			if err != nil {
				log.Println("Error fetching entries:", err)
				continue
			}

			for _, entry := range entries {
				
				//creating json to push into kafka topic
				data := model.KafkaPayload{
					Message: entry.Message,
					UserID:  entry.UserId,
				}
				
				jsonBytes, err := json.Marshal(data)
				if err != nil {
					log.Println("Error marshaling JSON:", err)
					return err
				}
				//err := publishToKafka(entry.Topic, entry.Message,entry.UserId)
				err = publishToKafka(entry.Topic, string(jsonBytes))
				if err != nil {
					log.Printf("Failed to publish entry %s: %v", entry.SchedulerEntryID, err)
					// optionally log failure
					continue
				}
				err = s.repo.MarkAsProcessed(ctx, entry.SchedulerEntryID)
				if err != nil {
					log.Printf("Failed to mark entry %s as processed: %v", entry.SchedulerEntryID, err)
				}
			}
		case <-ctx.Done():
			log.Println("Polling stopped")
			return nil
		}
	}
}




