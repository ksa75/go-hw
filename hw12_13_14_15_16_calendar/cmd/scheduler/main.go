package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mycalendar/internal/config"
	"mycalendar/internal/mq"
	"mycalendar/internal/scheduler"
	"mycalendar/internal/storage"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	conf, err := config.Read(configFile)
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	rmq, err := mq.NewRabbitMQ(conf.Queue.URL)
	fmt.Println(conf.Queue.URL)
	if err != nil {
		return fmt.Errorf("MQ error: %w", err)
	}
	defer rmq.Close()

	var store storage.EventsStorage

	sqlStore := sqlstorage.New()
	if err := sqlStore.Connect(ctx, conf.PSQL.DSN); err != nil {
		return fmt.Errorf("DB connect failed: %w", err)
	}
	store = sqlStore

	s := scheduler.NewScheduler(store, rmq, conf.Queue.Name, conf.Scheduler.CleanupOlderThanDays)

	ticker := time.NewTicker(time.Duration(conf.Scheduler.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Run(ctx)
		case <-ctx.Done():
			log.Println("scheduler: context cancelled, exiting")
			return nil
		}
	}
}
