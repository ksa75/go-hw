package internalhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"mycalendar/internal/app"
	internalhttp "mycalendar/internal/server/http"
	memorystorage "mycalendar/internal/storage/memory"
)

type testLogger struct{}

func (testLogger) Printf(_ string, _ ...any) {}
func (testLogger) Info(_ string)             {}
func (testLogger) Error(_ string)            {}

func TestIntegration_CreateGetDeleteEvent(t *testing.T) {
	// Создаём реальные компоненты
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	server := internalhttp.NewServer(logger, appInstance, "localhost", "8080")

	// --- POST /events: создаём событие
	startAt := time.Now().UTC().Truncate(time.Second)
	event := map[string]any{
		"userId":       "user1",
		"title":        "Integration test event",
		"description":  "desc",
		"startAt":      startAt.Format(time.RFC3339),
		"duration":     "1h",
		"noticeBefore": "10m",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// --- GET /events: проверяем, что событие есть
	req = httptest.NewRequest(http.MethodGet, "/events", nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var events []map[string]any
	err = json.NewDecoder(rec.Body).Decode(&events)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "Integration test event", events[0]["Title"])

	// --- DELETE /events: удаляем событие
	url := "/events?userId=user1&startAt=" + startAt.Format(time.RFC3339)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)

	// --- GET /events: проверяем, что событий нет
	req = httptest.NewRequest(http.MethodGet, "/events", nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	err = json.NewDecoder(rec.Body).Decode(&events)
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestIntegration_UpdateEvent(t *testing.T) {
	// Setup
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	server := internalhttp.NewServer(logger, appInstance, "localhost", "8080")
	startAt := time.Now().UTC().Truncate(time.Second)

	// Step 1: Create event
	createEvent := map[string]any{
		"userId":       "user1",
		"title":        "Original Title",
		"description":  "Original description",
		"startAt":      startAt.Format(time.RFC3339),
		"duration":     "1h",
		"noticeBefore": "10m",
	}
	body, err := json.Marshal(createEvent)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Step 2: Update event
	updateEvent := map[string]any{
		"userId":       "user1",
		"title":        "Updated Title",
		"description":  "Updated description",
		"startAt":      startAt.Format(time.RFC3339),
		"duration":     "2h",
		"noticeBefore": "20m",
	}
	body, err = json.Marshal(updateEvent)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Step 3: Fetch events and verify update
	req = httptest.NewRequest(http.MethodGet, "/events", nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var events []map[string]any
	err = json.NewDecoder(rec.Body).Decode(&events)
	require.NoError(t, err)
	require.Len(t, events, 1)

	updated := events[0]
	require.Equal(t, "Updated Title", updated["Title"])
	require.Equal(t, "Updated description", updated["Description"])
	require.Equal(t, "2h", updated["Duration"])
	require.Equal(t, "20m", updated["NoticeBefore"])

	// Cleanup
	url := "/events?userId=user1&startAt=" + startAt.Format(time.RFC3339)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestIntegration_GetEventsByDayWeekMonth(t *testing.T) {
	// Setup components
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	server := internalhttp.NewServer(logger, appInstance, "localhost", "8080")

	startAt := time.Now().UTC().Truncate(time.Second)

	// Create event
	event := map[string]any{
		"userId":       "user1",
		"title":        "Event for range testing",
		"description":  "Range test",
		"startAt":      startAt.Format(time.RFC3339),
		"duration":     "1h",
		"noticeBefore": "5m",
	}
	body, err := json.Marshal(event)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Helper to test a given path
	testRangeEndpoint := func(path string) {
		url := path + "?date=" + startAt.Format("2006-01-02")
		req := httptest.NewRequest(http.MethodGet, url, nil)
		rec := httptest.NewRecorder()
		server.Handler().ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var events []map[string]any
		err := json.NewDecoder(rec.Body).Decode(&events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "Event for range testing", events[0]["Title"])
	}

	// Test day/week/month endpoints
	testRangeEndpoint("/events/day")
	testRangeEndpoint("/events/week")
	testRangeEndpoint("/events/month")

	// Cleanup
	url := "/events?userId=user1&startAt=" + startAt.Format(time.RFC3339)
	req = httptest.NewRequest(http.MethodDelete, url, nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}
