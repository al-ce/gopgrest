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

// ServeHTTP routes the request by method and path, where the path begins with
// an existing table name
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exists, err := h.tableExists(r)
	switch {
	case !exists || err != nil:
		log.Println(r.URL.Path, "not found")
		NotFoundHandler(w, r)
	case r.Method == http.MethodGet && ReListRequest.MatchString(r.URL.Path):
		h.Read(w, r)
	case r.Method == http.MethodPost:
		h.Create(w, r)
	case r.Method == http.MethodDelete && ReRequestWithId.MatchString(r.URL.Path):
		h.Delete(w, r)
	case r.Method == http.MethodPut && ReRequestWithId.MatchString(r.URL.Path):
		h.Update(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// Update updates a row in the table by id
func (h *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Get id
	matches := ReRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}
	setID := matches[1]

	// Decode request body into map to dynamically update row
	var updateData map[string]any
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Update row with request data
	if err := h.service.UpdateRow(table, setID, updateData); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// Delete adds removes a row from a table by id
func (h *APIHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Get id
	matches := ReRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}
	setID := matches[1]

	// Delete row by id
	if err := h.service.DeleteRow(table, setID); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// Create adds a row to a table
func (h *APIHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Decode request
	var data *map[string]any
	err = json.NewDecoder(r.Body).Decode(&data)
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

// Read gets rows from a table in the database, optionally filtering by query
// params
func (h *APIHandler) Read(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Retrieve listQueryResults from database
	listQueryResults, err := h.service.ListRows(table, r.URL.Query())
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(listQueryResults)
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
func (h *APIHandler) tableExists(r *http.Request) (bool, error) {
	// Get tableName from URL path
	tableName, err := h.extractTableName(r)
	if err != nil {
		return false, err
	}
	table, err := h.repo.GetTable(tableName)
	return table != nil, err
}

// extractTableName gets the table name from the URL
func (h *APIHandler) extractTableName(r *http.Request) (string, error) {
	matches := ReTable.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract table name from %s", r.URL.Path)
	}
	return matches[1], nil
}
