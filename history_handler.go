package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const fp = "history.json"

type HistoryHandler struct{}

// ServeHTTP routes the request by method and path
func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/":
		h.GetHistory(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// GetHistory reads a history file as defined by the subpath and responds with JSON
func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	// Read history file contents
	history, err := os.ReadFile(fp)
	if err != nil {
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(history)
}
