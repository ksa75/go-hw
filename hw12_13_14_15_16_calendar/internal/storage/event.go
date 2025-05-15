package storage

import (
	"context"
	"time"
)

type EventsStorage interface {
	AddEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	DeleteEvent(ctx context.Context, userID string, start time.Time) error
	GetEvents(ctx context.Context) ([]Event, error)
	GetEventsByDay(ctx context.Context, date time.Time) ([]Event, error)
	GetEventsByWeek(ctx context.Context, date time.Time) ([]Event, error)
	GetEventsByMonth(ctx context.Context, date time.Time) ([]Event, error)
}

type BaseStorage interface {
	Connect(ctx context.Context, dsn string) error
	Migrate(ctx context.Context, migrate string) error
	Close() error
	EventsStorage
}

type Event struct {
	UserID        string
	Title         string
	Description   string
	StartDateTime time.Time
	Duration      string
	NoticeBefore  string
	CreatedAt     time.Time
}
