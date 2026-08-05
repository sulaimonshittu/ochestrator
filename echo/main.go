package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

type Message struct {
	Msg string
}

type Config struct {
	addr string
}

func home(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	m := Message{}
	err := d.Decode(&m)
	if err != nil {
		json.NewEncoder(w).Encode(errors.New("Unable to decode request body"))
		return
	}
	log.Printf("Received message: %v\n", m.Msg)

	json.NewEncoder(w).Encode(m)
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check called\n")
	w.Write([]byte("OK"))
}

func healthFail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check failed")
	w.WriteHeader(http.StatusInternalServerError)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func main() {
	r := chi.NewRouter()

	r.Post("/", home)
	r.Get("/health", health)
	r.Get("/healthfail", healthFail)

	cfg := Config{}
	flag.StringVar(&cfg.addr, "addr", ":8080", "Server listening address")
	flag.Parse()

	go func() {
		log.Println("Listening on http://localhost:7777")
		if err := http.ListenAndServe(cfg.addr, r); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
	<-c

	log.Println("Shutting down")
}
