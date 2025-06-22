package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"mycalendar/internal/storage"
)

type Logger interface {
	Printf(format string, v ...any)
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, uID, title, desc, dur string, noticeBefore int32, startAt time.Time) error
	UpdateEvent(ctx context.Context, uID, title, desc, dur string, noticeBefore int32, startAt time.Time) error
	DeleteEvent(ctx context.Context, userID string, start time.Time) error
	DeleteOldEvents(ctx context.Context, before time.Time) error
	GetEvents(ctx context.Context) ([]storage.Event, error)
	GetUpcomingEvents(ctx context.Context, from time.Time) ([]storage.Event, error)
	GetEventsByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	s := &Server{
		logger: logger,
		app:    app,
	}
	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger))

	// RESTful endpoints
	r.HandleFunc("/events", s.createEventHandler).Methods(http.MethodPost)
	r.HandleFunc("/events", s.updateEventHandler).Methods(http.MethodPut)
	r.HandleFunc("/events", s.deleteEventHandler).Methods(http.MethodDelete)
	r.HandleFunc("/events", s.getEventsHandler).Methods(http.MethodGet)

	// Чтение
	r.HandleFunc("/events/day", s.getEventsByDayHandler).Methods(http.MethodGet)
	r.HandleFunc("/events/week", s.getEventsByWeekHandler).Methods(http.MethodGet)
	r.HandleFunc("/events/month", s.getEventsByMonthHandler).Methods(http.MethodGet)

	// Простой hello
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	r.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	s.server = &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second, // защита от Slowloris
	}

	return s
}

func (s *Server) Handler() http.Handler {
	return s.server.Handler
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Printf("HTTP server error: %v", err)
		}
	}()
	s.logger.Printf("HTTP server listening on %s", s.server.Addr)

	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Printf("Server stopping...")
	return s.server.Shutdown(ctx)
}
