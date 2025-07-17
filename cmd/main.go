package main

import (
	"context"
	doconf "docshell/internal/v1/config"
	"docshell/internal/v1/docs/handlers"
	docshell "docshell/internal/v1/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TODO Move Graceful shutdown to function
func main() {
	cfg := doconf.Config.Service.Web

	// Opening server
	doc := docshell.New(cfg.Host, cfg.Port)

	// Apply middlewares
	doc.Use(docshell.LoggingMiddleware)
	doc.Use(docshell.RecoveryMiddleware)
	// Adding routes
	//
	doc.GetRouter().GET("/", handlers.GetAllDocuments)

	// Start server with goroutine
	go func() {
		if err := doc.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server %s:%d failed to start - %v.", cfg.Host, cfg.Port, err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT)
	defer stop()
	// Wait for context to be cancelled
	<-ctx.Done()

	// Create a new context with a timeout for the shutdown process
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	// Attempt to gracefully shut down server
	if err := doc.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed - %v.", err)
	}

	log.Println("Server gracefully stopped.")
}

func gracefulShutdown() {

}
