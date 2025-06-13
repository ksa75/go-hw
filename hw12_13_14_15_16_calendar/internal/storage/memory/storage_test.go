package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"mycalendar/internal/storage"
)

func TestStorage_AddEvent(t *testing.T) {
	mem := New()
	ctx := context.Background()

	event := storage.Event{
		UserID:        "u1",
		Title:         "meeting",
		StartDateTime: time.Date(2025, 5, 13, 15, 0, 0, 0, time.UTC),
	}

	err := mem.AddEvent(ctx, event)
	require.NoError(t, err)

	// попытка добавить в то же время
	err = mem.AddEvent(ctx, event)
	require.ErrorIs(t, err, storage.ErrDateBusy)

	events, err := mem.GetEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)
}

func TestStorage_UpdateAndDelete(t *testing.T) {
	mem := New()
	ctx := context.Background()

	start := time.Now().Truncate(time.Minute)

	event := storage.Event{
		UserID:        "u1",
		Title:         "call",
		StartDateTime: start,
	}

	require.NoError(t, mem.AddEvent(ctx, event))

	// update
	event.Title = "updated call"
	require.NoError(t, mem.UpdateEvent(ctx, event))

	evs, _ := mem.GetEvents(ctx)
	require.Equal(t, "updated call", evs[0].Title)

	// delete
	err := mem.DeleteEvent(ctx, "u1", start)
	require.NoError(t, err)

	evs, _ = mem.GetEvents(ctx)
	require.Len(t, evs, 0)
}

func TestStorage_DeleteOldEvents(t *testing.T) {
	mem := New()
	ctx := context.Background()

	old := time.Now().AddDate(-2, 0, 0)
	recent := time.Now().AddDate(0, 0, -1)

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "very old",
		StartDateTime: old,
	})

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "recent",
		StartDateTime: recent,
	})

	err := mem.DeleteOldEvents(ctx, time.Now().AddDate(-1, 0, 0)) // delete older than 1 year
	require.NoError(t, err)

	events, err := mem.GetEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "recent", events[0].Title)
}

func TestStorage_GetUpcomingEvents(t *testing.T) {
	mem := New()
	ctx := context.Background()

	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "past event",
		StartDateTime: past,
	})

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "future event",
		StartDateTime: future,
	})

	events, err := mem.GetUpcomingEvents(ctx, time.Now())
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "future event", events[0].Title)
}

func TestStorage_GetEventsByDay(t *testing.T) {
	mem := New()
	ctx := context.Background()

	targetDay := time.Date(2025, 6, 13, 10, 0, 0, 0, time.UTC)
	otherDay := time.Date(2025, 6, 14, 10, 0, 0, 0, time.UTC)

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "on day",
		StartDateTime: targetDay,
	})

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "another day",
		StartDateTime: otherDay,
	})

	events, err := mem.GetEventsByDay(ctx, targetDay)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "on day", events[0].Title)
}

func TestStorage_GetEventsByWeek(t *testing.T) {
	mem := New()
	ctx := context.Background()

	weekDate := time.Date(2025, 6, 12, 10, 0, 0, 0, time.UTC) // same week
	otherDate := time.Date(2025, 7, 1, 10, 0, 0, 0, time.UTC) // different week

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "this week",
		StartDateTime: weekDate,
	})

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "next month",
		StartDateTime: otherDate,
	})

	events, err := mem.GetEventsByWeek(ctx, weekDate)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "this week", events[0].Title)
}

func TestStorage_GetEventsByMonth(t *testing.T) {
	mem := New()
	ctx := context.Background()

	monthDate := time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC)
	otherDate := time.Date(2025, 7, 10, 10, 0, 0, 0, time.UTC)

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "June event",
		StartDateTime: monthDate,
	})

	mem.AddEvent(ctx, storage.Event{
		UserID:        "u1",
		Title:         "July event",
		StartDateTime: otherDate,
	})

	events, err := mem.GetEventsByMonth(ctx, monthDate)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "June event", events[0].Title)
}
