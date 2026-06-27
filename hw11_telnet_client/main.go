package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()
	_, err := fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "print failed: %v\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		for {
			err := client.Receive()
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Fprintf(os.Stderr, "connection closed\n")
				}
				stop()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := client.Send()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintf(os.Stderr, "connection closed\n")
				client.Close()
			}
			return
		}
	}

}
