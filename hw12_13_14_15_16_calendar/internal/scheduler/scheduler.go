package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"mycalendar/internal/mq"
	"mycalendar/internal/notifier"
	"mycalendar/internal/storage"
)

type Scheduler struct {
	storage   storage.EventsStorage
	publisher mq.Publisher
	topic     string
	days      int
}

func NewScheduler(s storage.EventsStorage, p mq.Publisher, t string, days int) *Scheduler {
	return &Scheduler{storage: s, publisher: p, topic: t, days: days}
}

func (s *Scheduler) Run(ctx context.Context) {
	events, err := s.storage.GetUpcomingEvents(ctx, time.Now())
	if err != nil {
		log.Printf("failed to get events: %v", err)
		return
	}

	for _, e := range events {
		notif := notifier.Notification{
			EventID: e.EventID,
			Title:   e.Title,
			StartAt: e.StartDateTime,
			UserID:  e.UserID,
		}
		data, err := json.Marshal(notif)
		if err != nil {
			log.Printf("marshal error: %v", err)
			continue
		}
		err = s.publisher.Publish(s.topic, data)
		if err != nil {
			log.Printf("publish error: %v", err)
		}
	}

	log.Println("sending remainders to queue")
	if s.days == 0 {
		return
	}
	err = s.storage.DeleteOldEvents(ctx, time.Now().AddDate(0, 0, -s.days))
	log.Println("deleting old events")
	if err != nil {
		log.Printf("error deleting old events: %v", err)
	}
}
