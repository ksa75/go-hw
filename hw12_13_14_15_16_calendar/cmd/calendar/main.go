package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"mycalendar/api/calendarpb"
	"mycalendar/internal/app"
	"mycalendar/internal/config"
	"mycalendar/internal/logger"
	grpcserver "mycalendar/internal/server/grpc"
	internalhttp "mycalendar/internal/server/http"
	"mycalendar/internal/storage"
	memorystorage "mycalendar/internal/storage/memory"
	sqlstorage "mycalendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.Read(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	////////////////////////
	mylogger, err := logger.New(conf.Logger.Level, conf.Logger.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	mylogger.Debug("this is debug")
	mylogger.Info("running on port")
	mylogger.Error("something went wrong")
	////////////////////////

	var store storage.EventsStorage

	switch conf.Storage.Type {
	case "memory":
		store = memorystorage.New()
	case "sql":
		sqlStore := sqlstorage.New()
		if err := sqlStore.Connect(ctx, conf.PSQL.DSN); err != nil {
			return fmt.Errorf("DB connect failed: %w", err)
		}
		if err := sqlStore.Migrate(ctx, conf.PSQL.Migration); err != nil {
			return fmt.Errorf("cannot migrate: %w", err)
		}
		store = sqlStore

	default:
		return fmt.Errorf("unknown storage type: %s", conf.Storage.Type)
	}

	////////////////////////
	calendar, err := app.New(mylogger, store)
	if err != nil {
		return fmt.Errorf("cannot create app: %w", err)
	}

	if conf.Storage.Type == "memory" {
		then, _ := time.Parse("2006-01-02", "2025-07-21")
		then1, _ := time.Parse("2006-01-02", "2025-12-21")
		err = calendar.CreateEvent(ctx, "007", "test event1", "test event", "1h", 15, then1)
		if err != nil {
			mylogger.Printf("failed to create event: %v", err)
		}
		err = calendar.CreateEvent(ctx, "006", "test event2", "test event", "1h", 15, then)
		if err != nil {
			mylogger.Printf("failed to create event: %v", err)
		}
		err = calendar.CreateEvent(ctx, "006", "test event3", "test event", "1h", 15, then1)
		if err != nil {
			mylogger.Printf("failed to create event: %v", err)
		}
	}

	if err := calendar.Run(ctx); err != nil {
		return fmt.Errorf("cannot run app: %w", err)
	}

	////////////////////////
	srv := internalhttp.NewServer(mylogger, calendar, conf.HTTP.Host, conf.HTTP.Port)

	grpcServer := grpc.NewServer()
	calendarpb.RegisterCalendarServiceServer(grpcServer, grpcserver.NewServer(calendar))

	// Create context with cancel on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start gRPC server in a goroutine
	go func() {
		grpcAddr := net.JoinHostPort(conf.GRPC.Host, conf.GRPC.Port)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			mylogger.Printf("failed to listen for gRPC: %v", err)
			stop()
			return
		}
		mylogger.Printf("gRPC server listening on %s", grpcAddr)

		if err := grpcServer.Serve(lis); err != nil {
			mylogger.Printf("gRPC server error: %v", err)
			stop()
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		if err := srv.Start(ctx); err != nil {
			mylogger.Printf("HTTP server error: %v", err)
			stop()
		}
	}()

	// Wait for signal
	<-ctx.Done()
	mylogger.Printf("Shutdown signal received")

	// Gracefully stop servers
	grpcServer.GracefulStop()
	if err := srv.Stop(context.Background()); err != nil {
		mylogger.Printf("error shutting down HTTP server: %v", err)
	}

	return nil
}
