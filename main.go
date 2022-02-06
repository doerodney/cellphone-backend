package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
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

func newMultiplexer() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", RootHandler)
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
