package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	address := flag.Arg(0) + ":" + flag.Arg(1)
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Connected to %s\n", address)

	// Handle Ctrl+C (SIGINT)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	done := make(chan struct{})

	// Read from stdin and send to server
	go func() {
		err := client.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
		}
		fmt.Fprintln(os.Stderr, "EOF received (Ctrl+D). Closing connection.")
		done <- struct{}{}
	}()

	// Read from server and write to stdout
	go func() {
		err := client.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
		}
		fmt.Fprintln(os.Stderr, "Connection closed by remote host.")
		done <- struct{}{}
	}()

	// Wait for either Ctrl+C or EOF/socket close
	select {
	case <-sigCh:
		fmt.Fprintln(os.Stderr, "\nReceived SIGINT, exiting.")
	case <-done:
	}

	client.Close()
}
