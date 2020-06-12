package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting greeter service ...")

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", HealthzHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(httpServer.ListenAndServe())
}

func HealthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
