package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host string
	Port string `required:"true"`

	Endpoint string `default:"/greet"`

	ShutdownTimeout time.Duration `default:"5s"`

	Greeting string `required:"true"`
}

func main() {
	log.Println("Starting greeter service ...")

	// Load the service configuration from environment variables
	var c Config
	if err := envconfig.Process("greeter", &c); err != nil {
		log.Fatalf("Failed to process config from environment variables: %s", err)
	}

	// Configure the HTTP multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc(c.Endpoint, NewGreeter(c.Greeting).Handler)

	// Configure the HTTP server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: mux,
	}

	// Start listening from connections and serve traffic
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error shutting down server: %s", err)
		}
	}()

	// Capture the system signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive it
	<-signalChan
	log.Println("Shutdown signal received, exiting...")

	// Configure a shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown the server: %s", err)
	}
}

func HealthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

type Greeter struct {
	message string
}

func NewGreeter(message string) *Greeter {
	return &Greeter{message: message}
}

func (g *Greeter) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(g.message))
}
