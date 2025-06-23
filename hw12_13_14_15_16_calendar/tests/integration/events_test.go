////go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"mycalendar/internal/storage"
	sqlstorage "mycalendar/internal/storage/sql"
)

type EventsIntegrationSuite struct {
	suite.Suite
	pool    *pgxpool.Pool
	storage storage.BaseStorage
}

func TestEventsIntegrationSuite(t *testing.T) {
	suite.Run(t, new(EventsIntegrationSuite))
}

func (s *EventsIntegrationSuite) SetupSuite() {
	s.storage = sqlstorage.New()
	err := s.storage.Connect(context.Background(), "postgres://postgres:postgres@postgres:5432/calendar")
	if err != nil {
		s.Fail("DB connect failed: %w", err)
	}

	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@postgres:5432/calendar")
	if err != nil {
		s.Fail("DB connect failed: %w", err)
	}
	s.pool = pool
}

func (s *EventsIntegrationSuite) TearDownSuite() {
	s.storage.Close()
}

func (s *EventsIntegrationSuite) TearDownTest() {
	_, err := s.pool.Exec(context.Background(), "TRUNCATE events CASCADE")
	if err != nil {
		log.Fatalf("failed to truncate table: %v", err)
	}
	fmt.Println("Table truncated")
}

func (s *EventsIntegrationSuite) TestAddAndGetEvent() {
	now := time.Now().Truncate(time.Second)

	event := storage.Event{
		UserID:        "user-1",
		Title:         "Meeting",
		Description:   "Team sync",
		StartDateTime: now.Add(1 * time.Hour),
		Duration:      "30m",
		NoticeBefore:  10,
		CreatedAt:     now,
	}

	err := s.storage.AddEvent(context.Background(), event)
	s.Require().NoError(err)

	dbEvent, _ := s.getDirectEvent(event.UserID, event.StartDateTime)

	s.Equal(event.UserID, dbEvent.UserID)
	s.Equal(event.Title, dbEvent.Title)
	s.Equal(event.Description, dbEvent.Description)
	s.Equal(event.StartDateTime.Format(time.RFC3339), dbEvent.StartDateTime.Format(time.RFC3339))
	s.Equal(event.Duration, dbEvent.Duration)
	s.Equal(event.NoticeBefore, dbEvent.NoticeBefore)
}

func (s *EventsIntegrationSuite) getDirectEvent(userID string, start time.Time) (storage.Event, bool) {
	query, args, err := sq.
		Select("id", "user_id", "title", "description", "start_date_time", "duration", "notice_before", "created_at").
		From("events").
		Where(sq.Eq{"user_id": userID, "start_date_time": start.Format(time.RFC3339)}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	s.Require().NoError(err)

	rows, err := s.pool.Query(context.Background(), query, args...)
	s.Require().NoError(err)
	defer rows.Close()

	var event storage.Event
	if rows.Next() {
		s.Require().NoError(rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.Title,
			&event.Description,
			&event.StartDateTime,
			&event.Duration,
			&event.NoticeBefore,
			&event.CreatedAt,
		))
		return event, true
	}

	// no rows found
	return storage.Event{}, false
}

func (s *EventsIntegrationSuite) TestUpdateEvent() {
	now := time.Now().Truncate(time.Second)

	event := storage.Event{
		UserID:        "user-1",
		Title:         "Initial",
		Description:   "Initial Desc",
		StartDateTime: now.Add(2 * time.Hour),
		Duration:      "1h",
		NoticeBefore:  15,
		CreatedAt:     time.Now(),
	}

	err := s.storage.AddEvent(context.Background(), event)
	s.Require().NoError(err)

	dbEvent, _ := s.getDirectEvent(event.UserID, event.StartDateTime)

	s.Equal(event.UserID, dbEvent.UserID)
	s.Equal(event.Title, dbEvent.Title)
	s.Equal(event.Description, dbEvent.Description)
	s.Equal(event.StartDateTime.Format(time.RFC3339), dbEvent.StartDateTime.Format(time.RFC3339))
	s.Equal(event.Duration, dbEvent.Duration)
	s.Equal(event.NoticeBefore, dbEvent.NoticeBefore)
}

func (s *EventsIntegrationSuite) TestDeleteEvent() {
	now := time.Now().Truncate(time.Second)
	ctx := context.Background()

	event := storage.Event{
		UserID:        "user-1",
		Title:         "To Delete",
		Description:   "To Delete",
		StartDateTime: now.Add(3 * time.Hour),
		Duration:      "45m",
		NoticeBefore:  5,
		CreatedAt:     time.Now(),
	}

	// Add the event
	err := s.storage.AddEvent(ctx, event)
	s.Require().NoError(err)

	// Delete the event
	err = s.storage.DeleteEvent(ctx, event.UserID, event.StartDateTime)
	s.Require().NoError(err)

	// Ensure the event has been deleted
	deletedEvent, _ := s.getDirectEvent(event.UserID, event.StartDateTime)
	s.Empty(deletedEvent.Title)
	s.Empty(deletedEvent.Description)
}

func (s *EventsIntegrationSuite) TestGetEventsByDay() {
	ctx := context.Background()
	today := time.Now().Truncate(24 * time.Hour)

	// Use a unique test ID to isolate test data
	testID := uuid.New().String()

	// Create 3 uniquely identifiable events
	expectedEvents := make([]storage.Event, 0, 3)
	for i := 0; i < 3; i++ {
		event := storage.Event{
			UserID:        "user-1",
			Title:         fmt.Sprintf("Daily Event %s #%d", testID, i),
			Description:   fmt.Sprintf("Test event for %s [%d]", testID, i),
			StartDateTime: today.Add(time.Duration(i) * time.Hour),
			Duration:      "30m",
			NoticeBefore:  10,
			CreatedAt:     time.Now(),
		}
		err := s.storage.AddEvent(ctx, event)
		s.Require().NoError(err)
		expectedEvents = append(expectedEvents, event)
	}

	// Get all events for the day
	allEvents, err := s.storage.GetEventsByDay(ctx, today)
	s.Require().NoError(err)

	// Filter only test-related events using the unique test ID
	var matchedEvents []storage.Event
	for _, e := range allEvents {
		if strings.Contains(e.Title, testID) {
			matchedEvents = append(matchedEvents, e)
		}
	}
	s.Require().Len(matchedEvents, 3)

	// Verify each expected event was found
	for _, expected := range expectedEvents {
		found := false
		for _, actual := range matchedEvents {
			if actual.UserID == expected.UserID &&
				actual.Title == expected.Title &&
				actual.Description == expected.Description &&
				actual.StartDateTime.Equal(expected.StartDateTime) &&
				actual.Duration == expected.Duration &&
				actual.NoticeBefore == expected.NoticeBefore {
				found = true
				break
			}
		}
		s.True(found, "Expected event %+v not found in matched results", expected)
	}
}

func (s *EventsIntegrationSuite) TestDeleteOldEvents() {
	now := time.Now().Truncate(time.Second)
	past := now.AddDate(0, 0, -365-1) // older than retention threshold 1y

	oldEvent := storage.Event{
		UserID:        "user-old",
		Title:         "Old Event",
		Description:   "Should be deleted",
		StartDateTime: past,
		Duration:      "1h",
		NoticeBefore:  5,
		CreatedAt:     past,
	}
	newEvent := storage.Event{
		UserID:        "user-new",
		Title:         "Recent Event",
		Description:   "Should remain",
		StartDateTime: now,
		Duration:      "1h",
		NoticeBefore:  5,
		CreatedAt:     now,
	}

	// Insert both events
	s.Require().NoError(s.storage.AddEvent(context.Background(), oldEvent))
	s.Require().NoError(s.storage.AddEvent(context.Background(), newEvent))

	// Run deletion with retention threshold
	cutoff := time.Now().AddDate(0, 0, -365)
	s.Require().NoError(s.storage.DeleteOldEvents(context.Background(), cutoff))

	// Validate old event is deleted
	_, found := s.getDirectEvent(oldEvent.UserID, oldEvent.StartDateTime)
	s.False(found, "Old event should have been deleted")

	// Validate recent event is still present
	remainingEvent, found := s.getDirectEvent(newEvent.UserID, newEvent.StartDateTime)
	s.True(found, "Recent event should still exist")
	s.Equal("Recent Event", remainingEvent.Title)
}
