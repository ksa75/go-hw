package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
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

func (a *App) CreateEvent(ctx context.Context, uID, title, desc, dur, noticeBefore string, startAt time.Time) error {
	now := time.Now()
	e := storage.Event{
		UserID:        uID,
		Title:         title,
		Description:   desc,
		StartDateTime: startAt,
		Duration:      dur,
		NoticeBefore:  noticeBefore,
		CreatedAt:     now,
	}
	return a.events.AddEvent(ctx, e)
}

func (a *App) UpdateEvent(ctx context.Context, uID, title, desc, dur, noticeBefore string, startAt time.Time) error {
	now := time.Now()
	e := storage.Event{
		UserID:        uID,
		Title:         title,
		Description:   desc,
		StartDateTime: startAt,
		Duration:      dur,
		NoticeBefore:  noticeBefore,
		CreatedAt:     now,
	}
	return a.events.UpdateEvent(ctx, e)
}

func (a *App) DeleteEvent(ctx context.Context, userID string, start time.Time) error {
	return a.events.DeleteEvent(ctx, userID, start)
}

func (a *App) GetEvents(ctx context.Context) ([]storage.Event, error) {
	return a.events.GetEvents(ctx)
}

func (a *App) GetEventsByDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.events.GetEventsByDay(ctx, date)
}

func (a *App) GetEventsByWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.events.GetEventsByWeek(ctx, date)
}

func (a *App) GetEventsByMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.events.GetEventsByMonth(ctx, date)
}

func (a *App) Run(ctx context.Context) error {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, nil)))
	events, err := a.events.GetEvents(ctx)
	if err != nil {
		return err
	}
	a.logger.Printf("events:")
	// log.Println("events:")
	for _, b := range events {
		a.logger.Printf("\t %+v", b)
	}

	return nil
}
