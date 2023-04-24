package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shop/backend-api/internal/handlers"
	"shop/backend-api/internal/repository"
	"shop/backend-api/internal/service"
	"shop/backend-api/pkg/config"
)

func main() {
	var configFile string = "config.json"

	cfg, err := config.Init(configFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("using %s file for set up", configFile)

	db, err := repository.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Panic(err.Error())
	}
	defer db.Close()
	log.Println("database is ready")

	dao := repository.NewDao(db)

	authService := service.NewAuthService(dao)
	productService := service.NewProductService(dao)

	app := handlers.NewClient(
		cfg,
		authService,
		productService,
	)

	server := app.Server()

	go func() {
		log.Printf("server started at %s\n", cfg.Address)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %s\n", err)
	} else {
		log.Println("server gracefully stoped")
	}
}
