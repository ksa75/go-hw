package internalhttp

import (
	"context"
	// "fmt"
	"net"
	"net/http"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Printf(format string, v ...any)
	Info(msg string)
	Error(msg string)
}

type Application interface {
	// Заглушка для будущей логики
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	mux := http.NewServeMux()

	// "/" и "/hello" возвращают одно и то же
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}

	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/hello", helloHandler)

	// функциональный конвейер
	handler := loggingMiddleware(logger)(mux)

	return &Server{
		logger: logger,
		app:    app,
		server: &http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: handler,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()
	s.logger.Printf("Server started at %s", s.server.Addr)

	<-ctx.Done()
	return s.Stop(context.Background())
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Printf("Server stopping...")
	return s.server.Shutdown(ctx)
}
