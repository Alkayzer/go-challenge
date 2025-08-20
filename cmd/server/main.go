package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"events-service/internal/config"
	"events-service/internal/database"
	"events-service/internal/handlers"
	"events-service/internal/middleware"
	"events-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	repo, err := database.NewPostgresRepository(ctx, cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to create database repository: %v", err)
	}
	defer repo.Close()

	eventService := service.NewEventService(repo)

	eventHandler := handlers.NewEventHandler(eventService)

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.JSONMiddleware)

	router.HandleFunc("/health", eventHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/events", eventHandler.CreateEvent).Methods("POST")
	router.HandleFunc("/events", eventHandler.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", eventHandler.GetEventByID).Methods("GET")

	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.GetServerAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
