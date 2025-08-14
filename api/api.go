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
	log.Println(r.Method, r.URL, r.RemoteAddr)

	urlString := r.URL.String()
	isRequestWithID := ReRequestWithId.MatchString(urlString)
	isListRequest := ReListRequest.MatchString(urlString)

	exists, err := h.tableExists(r)
	switch {
	case r.Method == http.MethodGet && urlString == "/":
		h.ShowTables(w, r)
	case !exists || err != nil:
		log.Println(urlString, "not found")
		NotFoundHandler(w, r)
	case r.Method == http.MethodGet && isRequestWithID:
		h.Pick(w, r)
	case r.Method == http.MethodGet && isListRequest:
		h.List(w, r)
	case r.Method == http.MethodPost:
		h.Insert(w, r)
	case r.Method == http.MethodDelete && isRequestWithID:
		h.Delete(w, r)
	case r.Method == http.MethodPut && isRequestWithID:
		h.Update(w, r)
	default:
		NotFoundHandler(w, r)
	}
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

// Update updates a row in the table by id
func (h *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	updateQueryResult, err := h.Service.UpdateRow(table, id, updateData)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Encode to JSON
	jsonData, err := json.Marshal(updateQueryResult)
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

// Delete adds removes a row from a table by id
func (h *APIHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	rowsAffected, err := h.Service.DeleteRow(table, id)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "deleted %d rows from table %s\n", rowsAffected, table)
}

// Insert adds a row to a table
func (h *APIHandler) Insert(w http.ResponseWriter, r *http.Request) {
	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Store body for potential multiple reads
	bodyBytes, _ := io.ReadAll(r.Body)
	// Set a fresh ReadCloser with the body bytes
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	r.Body.Close()

	var data *[]types.RowData
	// Try decoding an array of JSON objects
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {

		// The request may not have been an array of JSON objects
		// Try decoding a single JSON object

		// Set a fresh ReadCloser with the body bytes
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		r.Body.Close()
		var singleRow *types.RowData

		if err = json.NewDecoder(r.Body).Decode(&singleRow); err != nil {

			// If we fail again, give up
			log.Println("Decode second attempt err", err)
			InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
			return

		} else {
			// If it was a single object, assign it as the only item in the
			// data array
			data = &[]types.RowData{*singleRow}
		}
	}

	// Insert new rows into the database one at a time so we can validate data
	newIds := []int64{}
	for _, row := range *data {
		newRowId, err := h.Service.InsertRow(&row, table)
		newIds = append(newIds, newRowId)
		if err != nil {
			log.Println(err)
			InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
			return
		}
	}

	// Set response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "rows created in table %s: %v", table, newIds)
}

// Pick gets a single row from a table in the database by id
func (h *APIHandler) Pick(w http.ResponseWriter, r *http.Request) {
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
	// Get table from URL path
	table, err := h.extractTableName(r)
	if err != nil {
		InternalServerErrorHandler(w, r, err.Error())
	}

	// Retrieve listQueryResults from database
	listQueryResults, err := h.Service.ListRows(table, r.URL.String())
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
