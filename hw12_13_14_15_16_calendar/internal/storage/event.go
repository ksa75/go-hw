package storage

import (
	"context"
	"time"
)

type EventsStorage interface {
	GetEvents(ctx context.Context) ([]Event, error)
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
