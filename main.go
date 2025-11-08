package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var serverIsBusy bool = false

type apiConfig struct {
	fileserverHits atomic.Int32
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {

	if serverIsBusy {
		http.Error(w, "Service temporarily anavailable", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewServer() *http.Server {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler)

	fs := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fs))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func main() {

	server := NewServer()

	fmt.Println("Server running on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Println("Error: ", err)
	}
}
