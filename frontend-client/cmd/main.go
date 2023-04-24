package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shop/frontend-client/internal"
)

func main() {
	server := internal.Server()

	go func() {
		fmt.Println("server started at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown: %s\n", err)
	} else {
		fmt.Println("server gracefully stoped")
	}
}
