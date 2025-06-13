package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mycalendar/internal/config"
	"mycalendar/internal/mq"
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
	if err != nil {
		return fmt.Errorf("can't connect to MQ: %w", err)
	}
	defer rmq.Close()

	handler := func(msg []byte) {
		fmt.Printf("📧 Уведомление: %s\n", string(msg))
	}

	// Start consuming in the background
	go func() {
		if err := rmq.Consume(conf.Queue.Name, handler); err != nil {
			log.Printf("consume failed: %v", err)
			cancel() // cancel context to exit main
		}
	}()

	// Block until context is done (e.g., signal received)
	<-ctx.Done()
	log.Println("shutting down gracefully")
	return nil
}
