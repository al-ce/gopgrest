package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"ftrack/repository"
	"ftrack/service"
)

type APIHandler struct {
	service service.Service
	repo    repository.Repository
}

var (
	ReRequestWithId = regexp.MustCompile(`^/\w+/([0-9]+)$`)
	ReListRequest   = regexp.MustCompile(`^/\w+(\?.*)?$`)
	ReTable         = regexp.MustCompile(`^/(\w+).*$`)
)

func NewAPIHandler(db *sql.DB) APIHandler {
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	return APIHandler{
		service: service,
		repo:    repo,
	}
}

// ServeHTTP routes the request by method and path
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch {
	case !h.tableExists(r):
		NotFoundHandler(w, r)
	case r.Method == http.MethodGet && ReListRequest.MatchString(r.URL.Path):
		h.ListSets(w, r)
	case r.Method == http.MethodPost:
		h.InsertRow(w, r)
	case r.Method == http.MethodDelete && ReRequestWithId.MatchString(r.URL.Path):
		h.DeleteSet(w, r)
	case r.Method == http.MethodPut && ReRequestWithId.MatchString(r.URL.Path):
		h.UpdateSet(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// UpdateSet adds an exercise set to the database by id
func (h *APIHandler) UpdateSet(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get id
	matches := ReRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}
	setID := matches[1]

	// Decode request body into map to dynamically update row
	var updateData map[string]any
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Update row with request data
	if err := h.service.UpdateSet(setID, updateData); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// DeleteSet adds an exercise set to the database by id
func (h *APIHandler) DeleteSet(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get id
	matches := ReRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}
	setID := matches[1]

	// Delete row by id
	if err := h.service.DeleteSet(setID); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// InsertRow adds an exercise set to the database
func (h *APIHandler) InsertRow(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table := r.URL.Path[1:]

	// Decode request
	var data *map[string]any
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Insert new set into the database
	if err = h.service.InsertRow(data, table); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// ListSets retrieves the exercise set history from the database, optionally
// filtering by query params
func (h *APIHandler) ListSets(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Retrieve sets from database
	sets, err := h.service.ListSets(r.URL.Query())
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(sets)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// tableExists checks if a resource references an existing table
func (h *APIHandler) tableExists(r *http.Request) bool {
	table := r.URL.Path[1:]
	return h.repo.TableExists(table)
}
