package main

import (
	"context"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"

	// "time"
	"fmt"
	"log"

	// "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	// memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"
	// "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	"mycalendar/internal/app"
	"mycalendar/internal/config"
	internalhttp "mycalendar/internal/server/http"
	sqlstorage "mycalendar/internal/storage/sql"
)

var configFile string

type dummyApp struct{}

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	////////////////////////
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
	////////////////////////
	// config := NewConfig()
	// logg := logger.New(config.Logger.Level)

	// storage := memorystorage.New()
	// calendar := app.New(logg, storage)

	// server := internalhttp.NewServer(logg, calendar)

	// ctx, cancel := signal.NotifyContext(context.Background(),
	// 	syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	// go func() {
	// 	<-ctx.Done()

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()

	// 	if err := server.Stop(ctx); err != nil {
	// 		logg.Error("failed to stop http server: " + err.Error())
	// 	}
	// }()

	// logg.Info("calendar is running...")

	// if err := server.Start(ctx); err != nil {
	// 	logg.Error("failed to start http server: " + err.Error())
	// 	cancel()
	// 	os.Exit(1) //nolint:gocritic
	// }
}

func mainImpl() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.Read(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %v", err)
	}

	s := new(sqlstorage.Storage)
	if err := s.Connect(ctx, c.PSQL.DSN); err != nil {
		return fmt.Errorf("cannot connect to psql: %v", err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			log.Println("cannot close psql connection", err)
		}
	}()

	if err := s.Migrate(ctx, c.PSQL.Migration); err != nil {
		return fmt.Errorf("cannot migrate: %v", err)
	}
	////////////////////////
	calendar, err := app.New(s)
	if err != nil {
		return fmt.Errorf("cannot create app: %v", err)
	}

	if err := calendar.Run(ctx); err != nil {
		return fmt.Errorf("cannot run app: %v", err)
	}

	////////////////////////
	// Открываем файл для логирования
	logFile, err := os.OpenFile("logs/access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("cannot open log file: %v", err)
	}
	defer logFile.Close()

	// Лог в stdout и в файл
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", 0)

	srv := internalhttp.NewServer(logger, calendar, "0.0.0.0", "8080")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		logger.Printf("Server exited with error: %v", err)
	}

	return nil
}
