package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type CellPhone struct {
	Id          int    `json:"id"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	OS          string `json:"os"`
	ReleaseDate string `json:"releaseDate"`
	Image       string `json:"image"`
}

func GetIpPortText(ip string, port int) string {
	text := fmt.Sprintf("%s:%d", ip, port)
	return text
}

func newMultiplexer() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/health", HealthHandler).Methods("GET")
	m.HandleFunc("/utc", UtcHandler).Methods("GET")
	m.HandleFunc("/api/phones", GetCellPhonesHandler).Methods("GET")
	m.HandleFunc("/api/phones/{id:[0-9]+}", GetCellPhoneByIdHandler).Methods("GET")
	m.HandleFunc("/api/phones/make/{make:[a-zA-Z]+}", GetCellPhonesByMakeHandler).Methods("GET")
	m.HandleFunc("/api/phones/os/{os:[a-zA-Z]+}", GetCellPhonesByOsHandler).Methods("GET")
	m.HandleFunc("/api/phones", PostCellPhoneHandler).Methods("POST")
	return m
}

func main() {
	m := newMultiplexer()
	port := 4000
	ip := "127.0.0.1"
	addr := GetIpPortText(ip, port)
	srv := &http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	fmt.Printf("Server is listening on %s\n", addr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
