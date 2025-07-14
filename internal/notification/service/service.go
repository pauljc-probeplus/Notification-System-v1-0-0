package service

import (
	"context"
	"notification-system/internal/notification/model"
	"notification-system/internal/notification/repository"
	upref_model "notification-system/internal/userpreference/model"
	upref_svc "notification-system/internal/userpreference/service" //Alias shall be given if multiple services are being imported
	sched_model "notification-system/internal/scheduler/model"
	sched_svc "notification-system/internal/scheduler/service"
	

	"time"
	"fmt"

	
)

type NotificationService interface {
	CreateNotification(ctx context.Context, n *model.Notification) error
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

func (s *notificationService) handleScheduling(ctx context.Context, n *model.Notification, pref *upref_model.UserPreference) error {
	//Logic from docx
	if pref.Preferences.NoDisturb.Enabled {
		return nil // placeholder
	}else{
		
		switch n.Type{
		case "transactional":
			entry := &sched_model.SchedulerEntry{
				NotificationID: n.NotificationID,
				UserId:         n.UserId,
				Message:        n.Message,
				SendAt:         n.SendAt,
				Status:         "pending",
			}
			for _, channel := range n.Channels {
				topic := "highPriority." + channel
				entry.Topic = topic// reset per channel

				if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
					return err // optionally log and continue instead of returning immediately
				}
			}
		case "system_alert":
			entry := &sched_model.SchedulerEntry{
				NotificationID: n.NotificationID,
				UserId:         n.UserId,
				Message:        n.Message,
				SendAt:         n.SendAt,
				Status:         "pending",
			}
			for _, channel := range pref.Preferences.Channels.SystemAlerts {
				topic := "highPriority." + channel
				entry.Topic = topic// reset per channel

				if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
					return err // optionally log and continue instead of returning immediately
				}
			}
		case "promotional":
			// TODO: Handle promotional-specific logic like limits and delivery window
			
			//layout := "2006-01-02T15:04:05"
			layout:=time.RFC3339
			sendAtTime, err := time.Parse(layout, n.SendAt)
			if err != nil {
				return fmt.Errorf("invalid SendAt format: %v", err)
			}

			createdTime, err := time.Parse(layout, n.CreatedDate)
			if err != nil {
				return fmt.Errorf("invalid CreatedDate format: %v", err)
			}
			if sendAtTime.Sub(createdTime) >= 24*time.Hour {
				// SendAt is 24+ hours after creation, push to DLQ
				return nil
			}else{
				entry := &sched_model.SchedulerEntry{
					NotificationID: n.NotificationID,
					UserId:         n.UserId,
					Message:        n.Message,
					Status:         "pending",
				}

				//Parse StartTime and EndTime (only time)
				timeLayout := "15:04"
				startT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.StartTime)
				if err != nil {
					return fmt.Errorf("invalid start time format: %v", err)
				}

				endT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.EndTime)
				if err != nil {
					return fmt.Errorf("invalid end time format: %v", err)
				}

				// Create full time.Time values for comparison on same day as SendAt
				deliveryStart := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
				startT.Hour(), startT.Minute(), 0, 0, sendAtTime.Location())

				deliveryEnd := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
				endT.Hour(), endT.Minute(), 0, 0, sendAtTime.Location())

				//Compare and adjust SendAt if necessary
				if sendAtTime.Before(deliveryStart) || sendAtTime.After(deliveryEnd) {
					// Outside preferred delivery window
					// Reschedule to next day at delivery start time
					nextDay := sendAtTime.AddDate(0, 0, 1)
					sendAtTime = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
						startT.Hour(), startT.Minute(), 0, 0, nextDay.Location())
				}
					
				//Final Step: Convert back to string (if needed)
				entry.SendAt = sendAtTime.Format(layout)

				for _, channel := range pref.Preferences.Channels.Promotional {
					topic := "lowPriority." + channel
					entry.Topic = topic// reset per channel
		
					if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
						return err // optionally log and continue instead of returning immediately
					}
				}
			}
		}
	}
	return nil
}


func (s *notificationService) CreateNotification(ctx context.Context, n *model.Notification) error {
	// TODO: Add business logic for DND, limits, etc later

	// Step 1: Persist notification
	if err := s.repo.Create(ctx, n); err != nil {
		return err
	}

	// Step 2: Fetch user preferences
	pref, err := s.userPrefService.GetUserPreference(ctx, n.UserId)
	if err != nil {
		return err // optionally return specific error if not found
	}

	//Logic from docx
	// if pref.Preferences.NoDisturb.Enabled {
	// 	return nil // placeholder
	// }else{
		
	// 	switch n.Type{
	// 	case "transactional":
	// 		entry := &sched_model.SchedulerEntry{
	// 			NotificationID: n.NotificationID,
	// 			UserId:         n.UserId,
	// 			Message:        n.Message,
	// 			SendAt:         n.SendAt,
	// 			Status:         "pending",
	// 		}
	// 		for _, channel := range pref.Preferences.Channels.Transactional {
	// 			topic := "highPriority." + channel
	// 			entry.Topic = topic// reset per channel
	
	// 			if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
	// 				return err // optionally log and continue instead of returning immediately
	// 			}
	// 		}
	// 	case "system_alert":
	// 		entry := &sched_model.SchedulerEntry{
	// 			NotificationID: n.NotificationID,
	// 			UserId:         n.UserId,
	// 			Message:        n.Message,
	// 			SendAt:         n.SendAt,
	// 			Status:         "pending",
	// 		}
	// 		for _, channel := range pref.Preferences.Channels.SystemAlerts {
	// 			topic := "highPriority." + channel
	// 			entry.Topic = topic// reset per channel
	
	// 			if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
	// 				return err // optionally log and continue instead of returning immediately
	// 			}
	// 		}
	// 	case "promotional":
	// 		// TODO: Handle promotional-specific logic like limits and delivery window
			
	// 		layout := "2006-01-02T15:04:05"
	// 		sendAtTime, err := time.Parse(layout, n.SendAt)
	// 		if err != nil {
	// 			return fmt.Errorf("invalid SendAt format: %v", err)
	// 		}

	// 		createdTime, err := time.Parse(layout, n.CreatedDate)
	// 		if err != nil {
	// 			return fmt.Errorf("invalid CreatedDate format: %v", err)
	// 		}
	// 		if sendAtTime.Sub(createdTime) >= 24*time.Hour {
	// 			// SendAt is 24+ hours after creation, push to DLQ
	// 			return nil
	// 		}else{
	// 			entry := &sched_model.SchedulerEntry{
	// 				NotificationID: n.NotificationID,
	// 				UserId:         n.UserId,
	// 				Message:        n.Message,
	// 				Status:         "pending",
	// 			}

	// 			//Parse StartTime and EndTime (only time)
	// 			timeLayout := "15:04"
	// 			startT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.StartTime)
	// 			if err != nil {
	// 				return fmt.Errorf("invalid start time format: %v", err)
	// 			}

	// 			endT, err := time.Parse(timeLayout, pref.Preferences.DeliveryTiming.EndTime)
	// 			if err != nil {
	// 				return fmt.Errorf("invalid end time format: %v", err)
	// 			}

	// 			// Create full time.Time values for comparison on same day as SendAt
	// 			deliveryStart := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
	// 			startT.Hour(), startT.Minute(), 0, 0, sendAtTime.Location())

	// 			deliveryEnd := time.Date(sendAtTime.Year(), sendAtTime.Month(), sendAtTime.Day(),
	// 			endT.Hour(), endT.Minute(), 0, 0, sendAtTime.Location())

	// 			//Compare and adjust SendAt if necessary
	// 			if sendAtTime.Before(deliveryStart) || sendAtTime.After(deliveryEnd) {
	// 				// Outside preferred delivery window
	// 				// Reschedule to next day at delivery start time
	// 				nextDay := sendAtTime.AddDate(0, 0, 1)
	// 				sendAtTime = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
	// 					startT.Hour(), startT.Minute(), 0, 0, nextDay.Location())
	// 			}
					
	// 			//Final Step: Convert back to string (if needed)
	// 			entry.SendAt = sendAtTime.Format(layout)

	// 			for _, channel := range pref.Preferences.Channels.Promotional {
	// 				topic := "lowPriority." + channel
	// 				entry.Topic = topic// reset per channel
		
	// 				if err := s.schedulerService.CreateSchedulerEntry(ctx, entry); err != nil {
	// 					return err // optionally log and continue instead of returning immediately
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// return nil
	
	// Step 3: Apply logic to schedule or delay
	return s.handleScheduling(ctx, n, pref)
}

