package storage

import "errors"

var (
	ErrDateBusy = errors.New("the selected datetime is already booked")
	ErrNotFound = errors.New("event not found")
)
