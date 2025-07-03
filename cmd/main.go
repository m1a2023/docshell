package main

import (
	"context"
	docshell "docshell/internal/server"
	"docshell/internal/storage"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server variables
var (
	PORT 		int
	HOST 		string
)

func main() {
	// Graceful shutdown 
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT)
	defer stop()

	// Opening server
	doc := docshell.New(HOST, PORT)
	// Apply middlewares
	doc.Use(docshell.LoggingMiddleware)
	doc.Use(docshell.RecoveryMiddleware)
	// Adding routes
	//

	// Start server with goroutine 
	go func() {
		if err := doc.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server %s:%d failed to start - %v.", HOST, PORT, err)
		}
	} ()
	// Wait for context to be cancelled 
	<- ctx.Done()

	// Create a new context with a timeout for the shutdown process
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancelShutdown()
	// Attempt to gracefully shut down server
	if err := doc.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed - %v.", err)
	}

	log.Println("Server gracefully stopped.")
}


func init() {
	// Setting up flags for using 
	setupFlags()
}

func setupFlags() {
	
	// Host and Port
	flag.StringVar(&HOST, "host", "0.0.0.0", "IP address to serve on")
	flag.StringVar(&HOST, "h", "0.0.0.0", "IP address shortcut")
	flag.IntVar(&PORT, "port", 8080, "Port to listen on")
	flag.IntVar(&PORT, "p", 8080, "Port shortcut")

	// Database Config
	flag.IntVar(&storage.MAX_IDLE_TIME, "max_idle_time", 5, "Max idle time for DB connections (in minutes)")
	flag.IntVar(&storage.MAX_CONNECTION_LIFE, "max_conn_life", 30, "Max connection lifetime (in minutes)")
	flag.IntVar(&storage.MAX_OPEN_CONNECTIONS, "max_open_conns", 10, "Max open connections to the DB")
	flag.IntVar(&storage.MAX_IDLE_CONNECTIONS, "max_idle_conns", 5, "Max idle connections in the pool")

	// Usage
	flag.Usage = printHelp
	flag.Parse()
	
	fmt.Printf("Database settings:\n")
	fmt.Printf("  Max Idle Time: %d min\n", storage.MAX_IDLE_TIME)
	fmt.Printf("  Max Connection Lifetime: %d min\n", storage.MAX_CONNECTION_LIFE)
	fmt.Printf("  Max Open Connections: %d\n", storage.MAX_OPEN_CONNECTIONS)
	fmt.Printf("  Max Idle Connections: %d\n", storage.MAX_IDLE_CONNECTIONS)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: cmd [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -h, --host string               IP address to serve on (default \"%s\")\n", HOST)
	fmt.Fprintf(os.Stderr, "  -p, --port string               Port to listen on (default \"%d\")\n", PORT)
	fmt.Fprintf(os.Stderr, "      --max_idle_time int         Max idle time for DB connections (minutes, default %d)\n", storage.MAX_IDLE_TIME)
	fmt.Fprintf(os.Stderr, "      --max_conn_life int         Max connection lifetime (minutes, default %d)\n", storage.MAX_CONNECTION_LIFE)
	fmt.Fprintf(os.Stderr, "      --max_open_conns int        Max open connections to DB (default %d)\n", storage.MAX_OPEN_CONNECTIONS)
	fmt.Fprintf(os.Stderr, "      --max_idle_conns int        Max idle connections in pool (default %d)\n", storage.MAX_IDLE_CONNECTIONS)
	fmt.Fprintf(os.Stderr, "  --help                          Show this help message\n")
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "cmd - Serve your app with configurable options\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -h, --host string\n")
	fmt.Fprintf(os.Stderr, "        IP address to serve on (default \"%s\").\n\n", HOST)

	fmt.Fprintf(os.Stderr, "  -p, --port string\n")
	fmt.Fprintf(os.Stderr, "        Port to listen on (default \"%s\").\n\n", PORT)

	fmt.Fprintf(os.Stderr, "      --max_idle_time int\n")
	fmt.Fprintf(os.Stderr, "        Max idle time for DB connections in minutes (default %d).\n\n", storage.MAX_IDLE_TIME)

	fmt.Fprintf(os.Stderr, "      --max_conn_life int\n")
	fmt.Fprintf(os.Stderr, "        Max lifetime for a DB connection in minutes (default %d).\n\n", storage.MAX_CONNECTION_LIFE)

	fmt.Fprintf(os.Stderr, "      --max_open_conns int\n")
	fmt.Fprintf(os.Stderr, "        Maximum number of open DB connections (default %d).\n\n", storage.MAX_OPEN_CONNECTIONS)

	fmt.Fprintf(os.Stderr, "      --max_idle_conns int\n")
	fmt.Fprintf(os.Stderr, "        Maximum number of idle DB connections (default %d).\n\n", storage.MAX_IDLE_CONNECTIONS)

	fmt.Fprintf(os.Stderr, "Example:\n")
	fmt.Fprintf(os.Stderr, "  cmd --host=0.0.0.0 --port=8080 --dsn=store.db --max_idle_time=5 --max_conn_life=30 --max_open_conns=10 --max_idle_conns=5\n")
}