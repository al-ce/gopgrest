package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"ftrack/models"
	"ftrack/repository"
	"ftrack/service"
)

type APIHandler struct {
	service service.Service
}

var ReRequestWithId = regexp.MustCompile(`^/sets/([0-9]+)$`)

func NewAPIHandler(db *sql.DB) APIHandler {
	sr := repository.NewRepository(db)
	service := service.NewService(sr)
	return APIHandler{
		service: service,
	}
}

// ServeHTTP routes the request by method and path
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/sets":
		h.ListSets(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/sets":
		h.CreateSet(w, r)
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
	log.Println(r.Method, r.URL.Path)

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
	log.Println(r.Method, r.URL.Path)

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

// CreateSet adds an exercise set to the database
func (h *APIHandler) CreateSet(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	// Decode request data
	var setData models.ExerciseSet
	err := json.NewDecoder(r.Body).Decode(&setData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Insert new set into the database
	if err = h.service.CreateSet(setData); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// ListSets retrieves the exercise set history from the database
func (h *APIHandler) ListSets(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	// Retrieve sets from database
	sets, err := h.service.ListSets()
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
