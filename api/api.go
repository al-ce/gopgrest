package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"gopgrest/repository"
	"gopgrest/service"
	"gopgrest/types"
)

type APIHandler struct {
	Service service.Service
	Repo    repository.Repository
}

type headers map[string]string

var (
	reRequestWithId     = regexp.MustCompile(`^/(\w+)/([0-9]+)$`)
	reRequestWithParams = regexp.MustCompile(`^/(\w+)(\?.*)?$`)
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

	if r.Method == http.MethodGet && r.URL.Path == "/" {
		h.showTables(w)
		return
	}

	// Standardize URL
	err := coerceURLToQueryParams(r)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}
	r.URL, err = url.Parse(decodeURL(r.URL.String()))
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	// Route request
	switch r.Method {
	case http.MethodGet:
		h.getRows(w, r)
	case http.MethodPost:
		h.insertRows(w, r)
	case http.MethodDelete:
		h.deleteRows(w, r)
	case http.MethodPut:
		h.updateRows(w, r)
	default:
		notFoundHandler(w)
	}
}

func (h *APIHandler) getRows(w http.ResponseWriter, r *http.Request) {
	table, err := parseOptionalParamsRequest(r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
	}

	// Retrieve gotRows from database
	gotRows, err := h.Service.GetRowsByRSQL(table, r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(gotRows)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	headers := headers{"Content-Type": "application/json"}
	writeResponse(w, http.StatusOK, headers, jsonData)
}

// insertRows adds a row to a table
func (h *APIHandler) insertRows(w http.ResponseWriter, r *http.Request) {
	table, err := parseOptionalParamsRequest(r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
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

			// If we fail again, it was malformed JSON/JSON array
			writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
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
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	// Set response
	data := fmt.Appendf(nil, "rows created in table %s: %v", table, newIds)
	writeResponse(w, http.StatusOK, nil, data)
}

func (h *APIHandler) updateRows(w http.ResponseWriter, r *http.Request) {
	tableName, err := parseOptionalParamsRequest(r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
	}

	// Decode request body into map to dynamically update row
	var updateData *types.RowData
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
		return
	}

	// Update row with request data
	rowsAffected, err := h.Service.UpdateRowsByRSQL(tableName, r.URL.String(), updateData)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	data := fmt.Appendf(nil, "rows updated in table %s: %d", tableName, rowsAffected)
	writeResponse(w, http.StatusOK, nil, data)
}

func (h *APIHandler) deleteRows(w http.ResponseWriter, r *http.Request) {
	// Get table from URL path
	tableName, err := parseOptionalParamsRequest(r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, []byte(err.Error()))
	}

	// Delete rows by rsql conditions
	rowsAffected, err := h.Service.DeleteRowsByRSQL(tableName, r.URL.String())
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}

	data := fmt.Appendf(nil, "rows deleted in table %s: %d", tableName, rowsAffected)
	writeResponse(w, http.StatusOK, nil, data)
}

// showTables responds with a JSON object of the tables, their column names,
// and column types
func (h *APIHandler) showTables(w http.ResponseWriter) {
	jsonData, err := json.Marshal(h.Repo.TablesRepr)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, []byte(err.Error()))
		return
	}
	headers := headers{"Content-Type": "application/json"}
	writeResponse(w, http.StatusOK, headers, jsonData)
}

// notFoundHandler responds with a 404 status and an error message
func notFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

// coerceURLToQueryParams takes a url with an ID resource, e.g. `/authors/1`,
// and converts it to a url with optional query params, e.g. `/authors?id==1`.
// If the URL didn't have an ID resource, returns the URL as is.
func coerceURLToQueryParams(r *http.Request) error {
	// If it's already in an optional-params format, return as is
	if reRequestWithParams.MatchString(r.URL.String()) {
		return nil
	}
	matches := reRequestWithId.FindStringSubmatch(r.URL.Path)
	if len(matches) < 3 {
		return fmt.Errorf("Could not parse for id and table: %s", r.URL.Path)
	}
	tableName := matches[1]
	rowID := matches[2]

	key := ""
	// Add rsql key for GET method
	if r.Method == http.MethodGet {
		key = "where="
	}

	// Convert URL to rsql format
	var err error
	r.URL, err = url.Parse(fmt.Sprintf("/%s?%sid==%s", tableName, key, rowID))
	if err != nil {
		return err
	}
	return nil
}

// parseOptionalParamsRequest gets the table name from a request with a url
// that does not contain an id resource and has optional query params, e.g.
// `/authors` or `/authors?select=surname`
func parseOptionalParamsRequest(url string) (string, error) {
	matches := reRequestWithParams.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract table name from %s", url)
	}
	return matches[1], nil
}

// writeResponse writes headers and data
func writeResponse(w http.ResponseWriter, statusCode int, headers headers, data []byte) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}
