package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"gopgrest/repository"
	"gopgrest/service"
	"gopgrest/types"
)

type APIHandler struct {
	Service service.Service
	Repo    repository.Repository
}

var (
	ReRequestWithId     = regexp.MustCompile(`^/(\w+)/([0-9]+)$`)
	ReRequestWithParams = regexp.MustCompile(`^/\w+(\?.*)?$`)
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
	log.Println(r.Method, r.URL, r.RemoteAddr)

	urlString := r.URL.String()
	isRequestWithID := ReRequestWithId.MatchString(urlString)
	isRequestWithParams := ReRequestWithParams.MatchString(urlString)

	exists, err := h.tableExists(r)
	switch {
	case r.Method == http.MethodGet && urlString == "/":
		h.ShowTables(w, r)
	case !exists || err != nil:
		log.Println(urlString, "not found")
		NotFoundHandler(w, r)
	case r.Method == http.MethodGet && isRequestWithID:
		h.GetRowByID(w, r)
	case r.Method == http.MethodGet && isRequestWithParams:
		h.GetRowsByRSQL(w, r)
	case r.Method == http.MethodPost:
		h.Insert(w, r)
	case r.Method == http.MethodDelete && isRequestWithID:
		h.DeleteRowByID(w, r)
	case r.Method == http.MethodDelete && isRequestWithParams:
		h.DeleteRowsByRSQL(w, r)
	case r.Method == http.MethodPut && isRequestWithID:
		h.UpdateRowByID(w, r)
	case r.Method == http.MethodPut && isRequestWithParams:
		h.UpdateRowByRSQL(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// GetRowByID gets a single row from a table in the database by id
func (h *APIHandler) GetRowByID(w http.ResponseWriter, r *http.Request) {
	tableName, rowID, err := parseByIDRequest(r.URL.Path)
	if err != nil {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}

	// Retrieve pickQueryResult from database
	pickQueryResult, err := h.Service.GetRowByID(tableName, rowID)
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

// GetRowsByRSQL gets rows from a table in the database with optional query
// params
func (h *APIHandler) GetRowsByRSQL(w http.ResponseWriter, r *http.Request) {
	table, err := h.parseOptionalParamsRequest(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Retrieve gotRows from database
	gotRows, err := h.Service.GetRowsByRSQL(table, r.URL.String())
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(gotRows)
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

// Insert adds a row to a table
func (h *APIHandler) Insert(w http.ResponseWriter, r *http.Request) {
	// Get table from URL path
	table, err := h.parseOptionalParamsRequest(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Store body for potential multiple reads
	bodyBytes, _ := io.ReadAll(r.Body)
	// Set a fresh ReadCloser with the body bytes
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	r.Body.Close()

	var newRows []types.RowData
	// Try decoding an array of JSON objects
	if err = json.NewDecoder(r.Body).Decode(&newRows); err != nil {

		// The request may not have been an array of JSON objects
		// Try decoding a single JSON object

		// Set a fresh ReadCloser with the body bytes
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		r.Body.Close()
		var singleRow types.RowData

		if err = json.NewDecoder(r.Body).Decode(&singleRow); err != nil {

			// If we fail again, give up
			log.Println("Decode second attempt err", err)
			InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
			return

		} else {
			// If it was a single object, assign it as the only item in the
			// data array
			newRows = []types.RowData{singleRow}
		}
	}

	// Insert new rows into the database
	newIds, err := h.Service.InsertRows(newRows, table)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "rows created in table %s: %v", table, newIds)
}

// UpdateRowByID updates a row in the table by id
func (h *APIHandler) UpdateRowByID(w http.ResponseWriter, r *http.Request) {
	tableName, rowID, err := parseByIDRequest(r.URL.Path)
	if err != nil {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}

	url := fmt.Sprintf("/%s?id==%s", tableName, rowID)
	h.updateRows(w, r, url)
}

// UpdateRowByRSQL updates any rows in the table matching the conditions in an
// RSQL query
func (h *APIHandler) UpdateRowByRSQL(w http.ResponseWriter, r *http.Request) {
	h.updateRows(w, r, r.URL.String())
}

func (h *APIHandler) updateRows(w http.ResponseWriter, r *http.Request, url string) {
	// Get table from URL path
	tableName, err := h.parseOptionalParamsRequest(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Decode request body into map to dynamically update row
	var updateData *types.RowData
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Update row with request data
	rowsAffected, err := h.Service.UpdateRowsByRSQL(tableName, url, updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "rows updated in table %s: %d", tableName, rowsAffected)
}

// DeleteRowByID adds removes a row from a table by id
func (h *APIHandler) DeleteRowByID(w http.ResponseWriter, r *http.Request) {
	// Get tableName and ID from URL path
	tableName, rowID, err := parseByIDRequest(r.URL.Path)
	if err != nil {
		InternalServerErrorHandler(w, r, "Could not find id match")
		return
	}

	url := fmt.Sprintf("/%s?id==%s", tableName, rowID)
	h.deleteRows(w, r, url)
}

// DeleteRowsByRSQL deletes any rows in the table matching the conditions in
// an RSQL query
func (h *APIHandler) DeleteRowsByRSQL(w http.ResponseWriter, r *http.Request) {
	h.deleteRows(w, r, r.URL.String())
}

func (h *APIHandler) deleteRows(w http.ResponseWriter, r *http.Request, url string) {
	// Get table from URL path
	tableName, err := h.parseOptionalParamsRequest(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Delete rows by rsql conditions
	rowsAffected, err := h.Service.DeleteRowsByRSQL(tableName, url)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "rows deleted in table %s: %d", tableName, rowsAffected)
}

// ShowTables responds with a JSON object of the tables, their column names,
// and column types
func (h *APIHandler) ShowTables(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(h.Repo.TablesRepr)
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
	tableName, err := h.parseOptionalParamsRequest(r)
	if err != nil {
		return false, err
	}
	table, err := h.Repo.GetTable(tableName)
	return table != nil, err
}

// parseByIDRequest gets a table name and an ID from a request with a url that
// contains an id resource after the table name, e.g. `/authors/1`
func parseByIDRequest(url string) (string, string, error) {
	matches := ReRequestWithId.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("Could not parse for id and table: %s", url)
	}
	tableName := matches[1]
	rowID := matches[2]
	return tableName, rowID, nil
}

// parseOptionalParamsRequest gets the table name from a request with a url
// that does not contain an id resource and has optional query params, e.g.
// `/authors` or `/authors?select=surname`
func (h *APIHandler) parseOptionalParamsRequest(r *http.Request) (string, error) {
	matches := ReRequestWithParams.FindStringSubmatch(r.URL.Path)
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
