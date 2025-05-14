package app

import (
	"context"
	"log"
	"time"

	"mycalendar/internal/storage"
)

type App struct {
	events storage.EventsStorage
	logger Logger
}

type Logger interface {
	Printf(format string, v ...any)
	Info(msg string)
	Error(msg string)
}

func New(logger Logger, events storage.EventsStorage) (*App, error) {
	return &App{
		events: events,
		logger: logger,
	}, nil
}

func (a *App) CreateEvent(ctx context.Context, userID, title string) error {
	now := time.Now()
	e := storage.Event{
		UserID:        userID,
		Title:         title,
		StartDateTime: now.Add(1 * time.Hour),
		Duration:      "1h",
		NoticeBefore:  "15m",
		CreatedAt:     now,
	}
	// a.logger.Printf("dfdsdsds %v", e)
	return a.events.AddEvent(ctx, e)
}

func (a *App) Run(ctx context.Context) error {
	events, err := a.events.GetEvents(ctx)
	if err != nil {
		return err
	}
	// a.logger.Printf("fdsfdsff %v", events)
	log.Println("events:")
	for _, b := range events {
		log.Printf("\t %+v", b)
	}

	return nil
}
