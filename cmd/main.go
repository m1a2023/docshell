package main

import (
	"context"
	doconf "docshell/internal/v1/config"
	"docshell/internal/v1/docs/handlers"
	"docshell/internal/v1/middleware/cors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := doconf.Config.Service.Web

	// Opening server
	doc := chi.NewRouter()

	// Apply middlewares
	doc.Use(middleware.Logger)
	doc.Use(cors.CORSMiddleware)
	doc.Use(middleware.Recoverer)

	// Adding routes
	doc.Route("/docs", func(r chi.Router) {
		r.Get("/", handlers.GetAllDocuments)
		r.Get("/id/{id}", handlers.GetDocumentById)

		r.Post("/", handlers.CreateDocument)

		// With query parameter 'path'
		r.Get("/download", handlers.DownloadDocument)
	})

	// Start server with goroutine
	go func() {
		adr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		msg := "Server[%s] successfuly started"
		log.Printf(msg, adr)
		if err := http.ListenAndServe(adr, doc); err != nil {
			msg := "Server[%s] fault with %v"
			log.Fatalf(msg, adr, err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT)
	defer stop()
	// Wait for context to be cancelled
	<-ctx.Done()
	log.Println("Server gracefully stopped.")
}
