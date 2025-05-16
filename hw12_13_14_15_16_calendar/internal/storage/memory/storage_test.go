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
