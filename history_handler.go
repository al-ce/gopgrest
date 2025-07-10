package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

var GetHistoryByName = regexp.MustCompile(`^/history/([a-z0-9]+)$`)

type HistoryHandler struct{}

// ServeHTTP routes the request by method and path
func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && GetHistoryByName.MatchString(r.URL.Path):
		h.GetHistory(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// GetHistory reads a history file as defined by the subpath and responds with JSON
func (h *HistoryHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	// Get resource path with regex
	matches := GetHistoryByName.FindStringSubmatch(r.URL.Path)

	// Expect full string + 1 match group e.g. [history/sample sample]
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "need resource path")
		return
	}

	// Read file contents from resource path
	fp := matches[1]
	history, err := os.ReadFile(fp + ".json")
	if err != nil {
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(history)
}
