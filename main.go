package main

import (
	"log"
	"net/http"
)

// Set represents a performed set
type Set struct {
	Name    string  `json:"name"`
	Date    string  `json:"date"`
	Output  float32 `json:"output"`
	Unit    string  `json:"unit"`
	Reps    int8    `json:"reps"`
	Set     int8    `json:"set"`
	Notes   string  `json:"notes"`
	Day     string  `json:"day"`
	Program string  `json:"program"`
	Tags    string  `json:"tags"`
}

func main() {
	// Create server and routes
	mux := http.NewServeMux()
	mux.Handle("/", &HistoryHandler{})

	// Run server
	log.Println("Listening on port 8090...")
	err := http.ListenAndServe(":8090", mux)
	log.Fatal(err)
}
