package service

import (
	"context"
	"notification-system/internal/notification/model"
	"notification-system/internal/notification/repository"
	upref_model "notification-system/internal/userpreference/model"
	upref_svc "notification-system/internal/userpreference/service" //Alias shall be given if multiple services are being imported
	sched_model "notification-system/internal/scheduler/model"
	sched_svc "notification-system/internal/scheduler/service"
	
	"github.com/google/uuid"

	"notification-system/internal/common/timeutil"

	"time"
	"fmt"
	"sync/atomic"

	
)

var schedulerCounter uint64 //  global atomic counter for schedulerEntryId


type NotificationService interface {
	CreateNotification(ctx context.Context, n *model.Notification) (bool,error)
}

type notificationService struct {
	repo repository.NotificationRepository
	userPrefService   upref_svc.UserPreferenceService
	schedulerService  sched_svc.SchedulerEntryService
	
}

func NewNotificationService(
	repo repository.NotificationRepository,
	userPrefService upref_svc.UserPreferenceService,
	schedulerService sched_svc.SchedulerEntryService) NotificationService {
	return &notificationService{
		repo: repo,
		userPrefService:  userPrefService,
		schedulerService: schedulerService,
	}
}

func (s *notificationService) handleScheduling(ctx context.Context, n *model.Notification, pref *upref_model.UserPreference) (bool,error) {
	//Logic from docx
	if pref.Preferences.NoDisturb.Enabled {
		s.logSchedulerFailure(ctx, n, fmt.Errorf("notification unscheduled as DoNotDisturb was set to true"))
		return false,nil // placeholder
	}else{
		
		switch n.Type{
		case "transactional":
			entry := &sched_model.SchedulerEntry{
				NotificationID: n.NotificationID,
				NotificationType: n.Type,
				UserId:         n.UserId,
				Message:        n.Message,
				SendAt:         n.SendAt,
				MaximumRetries: n.MaximumRetries,
				Status:         "pending",
				CreatedByName: n.CreatedByName,
				CreatedById: n.CreatedByID,
				ModifiedByName: n.ModifiedByName,
				ModifiedById: n.ModifiedByID,
				
			}
			for _, channel := range n.Channels {
				topic := "highPriority." + channel
				entry.Topic = topic// reset per channel

				entry.SchedulerEntryID = generateSchedulerEntryID() //setting schedulerEntryId 

				entry.CreatedDate = timeutil.NowISTFormatted()
				entry.ModifiedDate = timeutil.NowISTFormatted()


				if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
					s.logSchedulerFailure(ctx, n, err)
					return false,err // optionally log and continue instead of returning immediately
				}
			}
		case "system_alert":
			entry := &sched_model.SchedulerEntry{
				NotificationID: n.NotificationID,
				NotificationType: n.Type,
				UserId:         n.UserId,
				Message:        n.Message,
				SendAt:         n.SendAt,
				MaximumRetries: n.MaximumRetries,
				Status:         "pending",
				CreatedByName: n.CreatedByName,
				CreatedById: n.CreatedByID,
				ModifiedByName: n.ModifiedByName,
				ModifiedById: n.ModifiedByID,
			}
			for _, channel := range pref.Preferences.Channels.SystemAlerts {
				topic := "highPriority." + channel
				entry.Topic = topic// reset per channel

				entry.SchedulerEntryID = generateSchedulerEntryID() //setting schedulerEntryId 

				entry.CreatedDate = timeutil.NowISTFormatted()
				entry.ModifiedDate = timeutil.NowISTFormatted()


				if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
					s.logSchedulerFailure(ctx, n, err)
					return false,err // optionally log and continue instead of returning immediately
				}
			}
		case "promotional":
			layout := "2006-01-02T15:04:05"
			sendAtTime, err := time.Parse(layout, n.SendAt)
			if err != nil {
				return false,fmt.Errorf("invalid SendAt format: %v", err)
			}
		
			createdTime, err := time.Parse(layout, n.CreatedDate)
			if err != nil {
				return false,fmt.Errorf("invalid CreatedDate format: %v", err)
			}
		
			if sendAtTime.Sub(createdTime) >= 24*time.Hour {
				// Over 24 hours – drop it
				s.logSchedulerFailure(ctx, n, fmt.Errorf("SendAt is more than 24 hours after CreatedDate"))
				return false,nil
			}
		
			// Prepare base entry
			entry := &sched_model.SchedulerEntry{
				NotificationID: n.NotificationID,
				NotificationType: n.Type,
				UserId:         n.UserId,
				Message:        n.Message,
				SendAt:         n.SendAt,
				MaximumRetries: n.MaximumRetries,
				Status:         "pending",
				CreatedByName: n.CreatedByName,
				CreatedById: n.CreatedByID,
				ModifiedByName: n.ModifiedByName,
				ModifiedById: n.ModifiedByID,
				
			}
		
			// Parse start and end times (only time of day)
			timeLayout := "15:04"
			startT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.StartTime)
			if err != nil {
				return false,fmt.Errorf("invalid start time format: %v", err)
			}
		
			endT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.EndTime)
			if err != nil {
				return false,fmt.Errorf("invalid end time format: %v", err)
			}
		
			// Construct time ranges on same day
			deliveryStart := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
				startT.Hour(), startT.Minute(), 0, 0, sendAtTime.Location())
			deliveryEnd := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
				endT.Hour(), endT.Minute(), 0, 0, sendAtTime.Location())
		
			// Adjust SendAt if needed
			if sendAtTime.Before(deliveryStart) {
				sendAtTime = deliveryStart
			} else if sendAtTime.After(deliveryEnd) {
				// Drop notification (outside delivery window)
				s.logSchedulerFailure(ctx, n, fmt.Errorf("SendAt (%s) is after delivery window end (%s)", sendAtTime.Format(layout), deliveryEnd.Format(layout)))
				return false,nil
			}
		
			// Set final SendAt
			entry.SendAt = sendAtTime.Format(layout)

			// Step 1: Convert notification channels to a map
			requested := make(map[string]bool)
			for _, ch := range n.Channels {
				requested[ch] = true
			}
		
			for _, channel := range pref.Preferences.Channels.Promotional {
				if(requested[channel]){
					entry.Topic = "lowPriority." + channel

					entry.SchedulerEntryID = generateSchedulerEntryID() //setting schedulerEntryId 

					entry.CreatedDate = timeutil.NowISTFormatted()
					entry.ModifiedDate = timeutil.NowISTFormatted()

					if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
						return false,err
					}
				}
			}
		}
	}
	return true,nil
}

// function to log failed scheduling
func (s *notificationService) logSchedulerFailure(ctx context.Context, n *model.Notification, err error) {
	if err == nil {
		err = fmt.Errorf("unknown scheduler failure occurred")
	}
	log := &sched_model.FailureLog{
		LogID:          uuid.NewString(),
		NotificationID: n.NotificationID,
		UserID:         n.UserId,
		Type:           n.Type,
		Message:        n.Message,
		FailureReason:  err.Error(),
		Timestamp:      time.Now().Format("2006-01-02T15:04:05"),
	}

	// Log the failure (fire and forget)
	_ = s.schedulerService.LogFailure(ctx,log)
}




func (s *notificationService) CreateNotification(ctx context.Context, n *model.Notification) (bool,error) {
	// TODO: Add business logic for DND, limits, etc later

	// Step 1: Persist notification
	if err := s.repo.Create(ctx, n); err != nil {
		return false,err
	}

	// Step 2: Fetch user preferences
	pref, err := s.userPrefService.GetUserPreference(ctx, n.UserId)
	if err != nil {
		return false,err // optionally return specific error if not found
	}
	
	// Step 3: Apply logic to schedule or delay
	return s.handleScheduling(ctx, n, pref)
}

//Function to set schedulerEntryId
func generateSchedulerEntryID() string {
	counter := atomic.AddUint64(&schedulerCounter, 1)
	return fmt.Sprintf("schedentry-%03d", counter%1000)
}


