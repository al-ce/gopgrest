package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"ftrack/repository"
	"ftrack/service"
	"ftrack/types"
)

type APIHandler struct {
	Service service.Service
	Repo    repository.Repository
}

var (
	ReRequestWithId = regexp.MustCompile(`^/\w+/([0-9]+)$`)
	ReListRequest   = regexp.MustCompile(`^/\w+(\?.*)?$`)
	ReTable         = regexp.MustCompile(`^/(\w+).*$`)
)

func NewAPIHandler(db repository.QueryExecutor, tables []repository.Table) APIHandler {
	repo := repository.NewRepository(db, tables)
	service := service.NewService(repo)
	return APIHandler{
		Service: service,
		Repo:    repo,
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
	case r.Method == http.MethodGet && ReRequestWithId.MatchString(r.URL.Path):
		h.Pick(w, r)
	case r.Method == http.MethodGet && ReListRequest.MatchString(r.URL.Path):
		h.List(w, r)
	case r.Method == http.MethodPost:
		h.Insert(w, r)
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
	id := matches[1]

	// Decode request body into map to dynamically update row
	var updateData *types.RowData
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Update row with request data
	if err := h.Service.UpdateRow(table, id, updateData); err != nil {
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
	id := matches[1]

	// Delete row by id
	if err := h.Service.DeleteRow(table, id); err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "row %s deleted from table %s\n", id, table)
}

// Insert adds a row to a table
func (h *APIHandler) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Decode request
	var data *types.RowData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Insert new set into the database
	newRowId, err := h.Service.InsertRow(data, table)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "row %d created in table %s", newRowId, table)
}

// Pick gets a single row from a table in the database by id
func (h *APIHandler) Pick(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	matches := ReRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}
	rowID := matches[1]

	// Retrieve pickQueryResult from database
	pickQueryResult, err := h.Service.PickRow(table, rowID)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(pickQueryResult)
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

// List gets rows from a table in the database, optionally filtering by query
// params
func (h *APIHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path, r.RemoteAddr)

	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Retrieve listQueryResults from database
	listQueryResults, err := h.Service.ListRows(table, types.QueryFilter(r.URL.Query()))
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
	table, err := h.Repo.GetTable(tableName)
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

// InternalServerErrorHandler responds with a 500 status and an error message
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

// NotFoundHandler responds with a 404 status and an error message
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
