package main

import (
	"log"
	"net/http"
)

// Set represents a performed set
type Set struct {
	Name    string
	Date    string
	Output  float32
	Unit    string
	Reps    int8
	Set     int8
	Notes   string
	Day     string
	Program string
	Tags    string
}

func main() {
	// Create server and routes
	mux := http.NewServeMux()
	mux.Handle("/history/", &HistoryHandler{})

	// Run server
	log.Println("Listening on port 8090...")
	err := http.ListenAndServe(":8090", mux)
	log.Fatal(err)
}
