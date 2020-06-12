package main

import (
	"log"
	"net"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host string
	Port string `required:"true"`

	Greeting string `required:"true"`
}

func main() {
	log.Println("Starting greeter service ...")

	var c Config
	if err := envconfig.Process("greeter", &c); err != nil {
		log.Fatalf("Failed to process config from environment variables: %s", err)
	}

	greeter := NewGreeter(c.Greeting)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.HandleFunc("/greet", greeter.Handler)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: mux,
	}

	log.Fatal(httpServer.ListenAndServe())
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
