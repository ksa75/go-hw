package scheduler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"mycalendar/internal/scheduler"
	"mycalendar/internal/storage"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetUpcomingEvents(ctx context.Context, from time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, from)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockStorage) DeleteOldEvents(ctx context.Context, before time.Time) error {
	args := m.Called(ctx, before)
	return args.Error(0)
}

func (m *MockStorage) AddEvent(_ context.Context, _ storage.Event) error    { return nil }
func (m *MockStorage) UpdateEvent(_ context.Context, _ storage.Event) error { return nil }
func (m *MockStorage) DeleteEvent(_ context.Context, _ string, _ time.Time) error {
	return nil
}
func (m *MockStorage) GetEvents(_ context.Context) ([]storage.Event, error) { return nil, nil }
func (m *MockStorage) GetEventsByDay(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *MockStorage) GetEventsByWeek(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (m *MockStorage) GetEventsByMonth(_ context.Context, _ time.Time) ([]storage.Event, error) {
	return nil, nil
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(queueName string, body []byte) error {
	args := m.Called(queueName, body)
	return args.Error(0)
}

func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestScheduler_Run_Success(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	mockPublisher := new(MockPublisher)

	event := storage.Event{
		EventID:       1,
		Title:         "Test Event",
		StartDateTime: time.Now().Add(time.Hour),
		UserID:        "user1",
	}

	mockStorage.On("GetUpcomingEvents", mock.Anything, mock.Anything).Return([]storage.Event{event}, nil)
	mockPublisher.On("Publish", "reminders", mock.Anything).Return(nil)
	mockStorage.On("DeleteOldEvents", mock.Anything, mock.Anything).Return(nil)

	s := scheduler.NewScheduler(mockStorage, mockPublisher, "reminders", 30)
	s.Run(ctx)

	mockStorage.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestScheduler_Run_StorageError(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	mockPublisher := new(MockPublisher)

	mockStorage.On("GetUpcomingEvents", mock.Anything, mock.Anything).Return([]storage.Event{}, errors.New("db error"))

	s := scheduler.NewScheduler(mockStorage, mockPublisher, "reminders", 30)
	s.Run(ctx)

	mockStorage.AssertExpectations(t)
}

func TestScheduler_Run_PublishError(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	mockPublisher := new(MockPublisher)

	event := storage.Event{
		EventID:       1,
		Title:         "Test Event",
		StartDateTime: time.Now().Add(time.Hour),
		UserID:        "user1",
	}

	mockStorage.On("GetUpcomingEvents", mock.Anything, mock.Anything).Return([]storage.Event{event}, nil)
	mockPublisher.On("Publish", "reminders", mock.Anything).Return(errors.New("publish error"))
	mockStorage.On("DeleteOldEvents", mock.Anything, mock.Anything).Return(nil)

	s := scheduler.NewScheduler(mockStorage, mockPublisher, "reminders", 30)
	s.Run(ctx)

	mockPublisher.AssertExpectations(t)
}

func TestScheduler_Run_NoDeletionIfZeroRetention(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	mockPublisher := new(MockPublisher)

	event := storage.Event{
		EventID:       1,
		Title:         "Test Event",
		StartDateTime: time.Now().Add(time.Hour),
		UserID:        "user1",
	}

	mockStorage.On("GetUpcomingEvents", mock.Anything, mock.Anything).Return([]storage.Event{event}, nil)
	mockPublisher.On("Publish", "reminders", mock.Anything).Return(nil)

	s := scheduler.NewScheduler(mockStorage, mockPublisher, "reminders", 0)
	s.Run(ctx)

	mockStorage.AssertNotCalled(t, "DeleteOldEvents", mock.Anything, mock.Anything)
}
