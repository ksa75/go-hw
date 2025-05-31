package internalhttp_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mycalendar/internal/app"
	internalhttp "mycalendar/internal/server/http"
	memorystorage "mycalendar/internal/storage/memory"

	"github.com/stretchr/testify/require"
)

type testLogger struct{}

func (testLogger) Printf(format string, v ...any) {}
func (testLogger) Info(msg string)                {}
func (testLogger) Error(msg string)               {}

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
		"user_id":       "user1",
		"title":         "Integration test event",
		"description":   "desc",
		"start_at":      startAt.Format(time.RFC3339),
		"duration":      "1h",
		"notice_before": "10m",
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
	url := "/events?user_id=user1&start=" + startAt.Format(time.RFC3339)
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
