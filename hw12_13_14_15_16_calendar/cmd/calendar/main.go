package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"mycalendar/internal/app"
	"mycalendar/internal/config"
	"mycalendar/internal/logger"
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

	// Пример вызова
	// err = calendar.CreateEvent(ctx, "user123", "demo")
	// if err != nil {
	// 	mylogger.Printf("failed to create event: %v", err)
	// }

	if err := calendar.Run(ctx); err != nil {
		return fmt.Errorf("cannot run app: %w", err)
	}

	////////////////////////
	srv := internalhttp.NewServer(mylogger, calendar, conf.HTTP.Host, conf.HTTP.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}
