package memorystorage

import (
	"context"
	"sync"
	"time"

	"mycalendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string][]storage.Event // key = userID
}

func New() *Storage {
	return &Storage{
		events: make(map[string][]storage.Event),
	}
}

func (s *Storage) AddEvent(ctx context.Context, e storage.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	// Проверка на занятость времени
	for _, ev := range s.events[e.UserID] {
		if ev.StartDateTime.Equal(e.StartDateTime) {
			return storage.ErrDateBusy
		}
	}
	s.events[e.UserID] = append(s.events[e.UserID], e)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	evs := s.events[e.UserID]
	for i, ev := range evs {
		if ev.StartDateTime.Equal(e.StartDateTime) {
			s.events[e.UserID][i] = e
			return nil
		}
	}
	return storage.ErrNotFound
}

func (s *Storage) DeleteEvent(ctx context.Context, userID string, start time.Time) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	evs := s.events[userID]
	for i, ev := range evs {
		if ev.StartDateTime.Equal(start) {
			// replacing append(s[:i], s[i+1]...) by slices.Delete(s, i, i+1), added in go1.21
			s.events[userID] = append(evs[:i], evs[i+1:]...)
			return nil
		}
	}
	return storage.ErrNotFound
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []storage.Event
	for _, evs := range s.events {
		result = append(result, evs...)
	}
	return result, nil
}

func (s *Storage) GetEventsByDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []storage.Event
	for _, evs := range s.events {
		for _, e := range evs {
			if e.StartDateTime.Year() == day.Year() &&
				e.StartDateTime.YearDay() == day.YearDay() {
				result = append(result, e)
			}
		}
	}
	return result, nil
}
