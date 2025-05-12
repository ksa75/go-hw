package main

import (
	"context"
	"flag"
	// "os"
	// "os/signal"
	// "syscall"
	// "time"
	"fmt"
	"log"
	

	// "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	// "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	// internalhttp "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/http"
	// memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"

	"mycalendar/internal/app"
    "mycalendar/internal/config"
	"mycalendar/internal/storage/sql"
)

var configFile string

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

	c, err := config.Read("configs/local.toml")
	if err != nil {
		return fmt.Errorf("cannot read config: %v", err)
	}

	r := new(sqlstorage.Storage)
	if err := r.Connect(ctx, c.PSQL.DSN); err != nil {
		return fmt.Errorf("cannot connect to psql: %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Println("cannot close psql connection", err)
		}
	}()

	if err := r.Migrate(ctx, c.PSQL.Migration); err != nil {
		return fmt.Errorf("cannot migrate: %v", err)
	}

	a, err := app.New(r)
	if err != nil {
		return fmt.Errorf("cannot create app: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		return fmt.Errorf("cannot run app: %v", err)
	}

	return nil
}
