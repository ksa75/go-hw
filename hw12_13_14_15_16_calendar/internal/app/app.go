package app

import (
	"context"
	"log"
	"mycalendar/internal/storage"
)

type App struct {
	s storage.BaseStorage
}

// type Logger interface { // TODO
// }

// type Storage interface { // TODO
// }

// func New(logger Logger, storage Storage) *App {
// 	return &App{}
// }

// func (a *App) CreateEvent(ctx context.Context, id, title string) error {
// 	// TODO
// 	return nil
// 	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
// }

// TODO

func New(s storage.BaseStorage) (*App, error) {
	return &App{s: s}, nil
}

func (a *App) Run(ctx context.Context) error {
	events, err := a.s.GetEvents(ctx)
	if err != nil {
		return err
	}

	log.Println("events:")
	for _, b := range events {
		log.Printf("\t %+v", b)
	}

	return nil
}
