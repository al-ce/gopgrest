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
	case r.Method == http.MethodGet && r.URL.Path == "/":
		h.ListSets(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

// ListSets retrieves the exercise set history from the database
func (h *APIHandler) ListSets(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	const listStmt = `
		select * from exercise_sets
	`
	rows, err := h.db.Query(listStmt)
	if err != nil {
		InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
		return
	}
	defer rows.Close()

	var sets []ExerciseSet
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
			log.Println(err)
			InternalServerErrorHandler(w, r, fmt.Sprintf("%v", err))
			return
		}
		sets = append(sets, *set)
	}

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
