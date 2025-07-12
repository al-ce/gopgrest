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
	default:
		NotFoundHandler(w, r)
	}
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

	// Execute delete query
	const deleteStmt = `delete from exercise_sets where id = $1`

	result, err := h.db.Exec(deleteStmt, setID)
	if err != nil {
		log.Println(err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	// Set response
	w.WriteHeader(http.StatusOK)
}

// CreateSet adds an exercise set to the database
func (h *APIHandler) CreateSet(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	// Decode request data
	var setData ExerciseSet
	err := json.NewDecoder(r.Body).Decode(&setData)
	if err != nil {
		log.Println(err)
		return
	}

	// Execute create query
	const createStmnt = `
		insert into exercise_sets
		(
			name,
			performed_at,
			weight,
			unit,
			reps,
			set_count,
			notes,
			split_day,
			program,
			tags
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Set time performed if not set
	if setData.PerformedAt.IsZero() {
		setData.PerformedAt = time.Now()
	}
	setData.PerformedAt = setData.PerformedAt.Round(time.Second)

	result, err := h.db.Exec(createStmnt,
		setData.Name,
		setData.PerformedAt,
		setData.Weight,
		setData.Unit,
		setData.Reps,
		setData.SetCount,
		setData.Notes,
		setData.SplitDay,
		setData.Program,
		setData.Tags,
	)
	if err != nil {
		log.Println(err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	// Set response
	w.WriteHeader(http.StatusOK)
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
			&set.ID,
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
