package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIHandler struct {
	db *sql.DB
}

func NewAPIHandler(db *sql.DB) APIHandler {
	return APIHandler{
		db: db,
	}
}

// ServeHTTP routes the request by method and path
func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/sets":
		h.ListSets(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// ListSets retrieves the exercise set history from the database
func (h *APIHandler) ListSets(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	// Execute query
	const listStmt = `
		select * from exercise_sets
	`
	rows, err := h.db.Query(listStmt)
	if err != nil {
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}
	defer rows.Close()

	// Scan rows into struct slice
	sets, err := scanExerciseSetRows(rows)
	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
	}
	log.Println("ListSets results length:", len(sets))

	// Encode to JSON
	jsonData, err := json.Marshal(sets)
	if err != nil {
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func scanExerciseSetRows(rows *sql.Rows) ([]ExerciseSet, error) {
	sets := []ExerciseSet{}
	for rows.Next() {
		set := &ExerciseSet{}
		err := rows.Scan(
			&set.Name,
			&set.PerformedAt,
			&set.Weight,
			&set.Unit,
			&set.Reps,
			&set.SetCount,
			&set.Notes,
			&set.SplitDay,
			&set.Program,
			&set.Tags,
		)
		if err != nil {
			return sets, err
		}
		sets = append(sets, *set)
	}
	return sets, nil
}
