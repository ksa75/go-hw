package internalhttp

import (
	"encoding/json"
	"net/http"
	"time"
)

type createOrUpdateEventRequest struct {
	UserID       string    `json:"userId"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartAt      time.Time `json:"startAt"`
	Duration     string    `json:"duration"`
	NoticeBefore string    `json:"noticeBefore"`
}

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var rq createOrUpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := s.app.CreateEvent(r.Context(), rq.UserID, rq.Title, rq.Description, rq.Duration, rq.NoticeBefore, rq.StartAt)
	if err != nil {
		http.Error(w, "failed to create event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201 Created
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var rq createOrUpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := s.app.UpdateEvent(r.Context(), rq.UserID, rq.Title, rq.Description, rq.Duration, rq.NoticeBefore, rq.StartAt)
	if err != nil {
		http.Error(w, "failed to update event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	startStr := r.URL.Query().Get("startAt")

	if userID == "" || startStr == "" {
		http.Error(w, "missing required query parameters: userId and startAt", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "invalid time format for 'startAt', use RFC3339", http.StatusBadRequest)
		return
	}

	if err := s.app.DeleteEvent(r.Context(), userID, startTime); err != nil {
		http.Error(w, "failed to delete event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (s *Server) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := s.app.GetEvents(r.Context())
	if err != nil {
		http.Error(w, "failed to get events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

func (s *Server) getEventsByDayHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsByDay(r.Context(), date)
	if err != nil {
		http.Error(w, "failed to get events by day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

func (s *Server) getEventsByWeekHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsByWeek(r.Context(), date)
	if err != nil {
		http.Error(w, "failed to get events by week: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

func (s *Server) getEventsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	date, err := parseDateQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := s.app.GetEventsByMonth(r.Context(), date)
	if err != nil {
		http.Error(w, "failed to get events by month: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, events)
}

func parseDateQuery(r *http.Request) (time.Time, error) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		return time.Time{}, http.ErrMissingFile // reuse existing error
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
