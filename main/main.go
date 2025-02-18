package main

import (
	"campaign"
	"campaign/logger"
	"campaign/session"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

type campaignHandler struct {
	authService *campaign.AuthService
}

func main() {
	godotenv.Load()

	authService, err := campaign.NewAuthService()
	if err != nil {
		logger.Println(err)
		return
	}

	session, err := session.New()
	if err != nil {
		logger.Println(err)
		return
	}

	log.Println(session)

	handler := &campaignHandler{
		authService: authService,
	}

	mux := mux.NewRouter()

	mux.HandleFunc("/campaign", campaignView).Methods("GET")

	mux.HandleFunc("/login", loginView).Methods("GET")
	mux.HandleFunc("/login", handler.stepHandler([]campaign.Step{
		authService.Login,
	})).Methods("POST")

	mux.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("template"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in a goroutine so it doesn't block shutdown handling
	go func() {
		fmt.Println("Server is running on port " + server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt) // Capture Ctrl+C (SIGINT) or SIGTERM

	// Block until a signal is received
	<-stop
	fmt.Println("\nShutting down server...")

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error during shutdown:", err)
	}

	fmt.Println("Server stopped gracefully.")
}
