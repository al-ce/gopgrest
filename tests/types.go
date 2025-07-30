package tests

import (
	"database/sql"
	"time"

	"ftrack/repository"
	"ftrack/types"
)

type TestDB struct {
	DB     *sql.DB
	TX     *sql.Tx
	Tables []repository.Table
}

// ExerciseSet matches the exercise_set table in the test database so that rows
// can be scanned into fields of appropriate size
type ExerciseSet struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PerformedAt time.Time `json:"performed_at"`
	Weight      int       `json:"weight"`
	Unit        string    `json:"unit"`
	Reps        int       `json:"reps"`
	SetCount    int       `json:"set_count"`
	Notes       string    `json:"notes"`
	SplitDay    string    `json:"split_day"`
	Program     string    `json:"program"`
	Tags        string    `json:"tags"`
}

type TagMap map[string]string

type FilterTest struct {
	TestName  string
	Filters   types.QueryFilter
	RowCount  int
	ExpectErr any
}
