package main

import (
	"os"
	"os/signal"
	"syscall"
	"context"
	"log"
	"time"
)

var sigs = []os.Signal{os.Interrupt, syscall.SIGTERM}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), sigs...)
	defer stop()

	log.Println("Press Ctrl+c to quit...")
	select {
	case <-time.After(5*time.Second):
		log.Println("Time up, quit anyway...")
	case <-ctx.Done():
		log.Printf("Receive quit signal(%v), bye...\n", ctx.Err())
		stop()
	}
}
