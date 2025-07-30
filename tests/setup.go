package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"ftrack/repository"
	"ftrack/service"
	"ftrack/types"
)

var (
	host       = os.Getenv("HOST")
	port       = os.Getenv("TEST_DB_PORT")
	user       = os.Getenv("TEST_DB_USER")
	password   = os.Getenv("TEST_DB_PASS")
	testDbName = os.Getenv("TEST_DB_NAME")
)

const TABLE1 = "exercise_sets"

type TestDB struct {
	DB *sql.DB
	TX *sql.Tx
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

// SampleRows are used to populate the test database
var SampleRows = []types.RowDataMap{
	{
		"Name":   "deadlift",
		"Weight": 300,
	},
	{
		"Name":   "deadlift",
		"Weight": 200,
	},
	{
		"Name":   "deadlift",
		"Weight": 100,
	},
	{
		"Name":   "squat",
		"Weight": 300,
	},
	{
		"Name":   "squat",
		"Weight": 200,
	},
	{
		"Name":   "squat",
		"Weight": 100,
	},
	// Entries we will NOT filter for
	{
		"Name":   "bench press",
		"Weight": 300,
	},
}

// NewTestDB returns a test database
func NewTestDB(t *testing.T) *TestDB {
	testParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, testDbName,
	)

	db, err := sql.Open("postgres", testParams)
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Get tables from database here rather than from a Tx later, which won't
	// return the column names
	tables, err := repository.GetPublicTables(db)

	t.Cleanup(func() {
		db.Close()
	})

	return &TestDB{
		db,
		nil,
		tables,
	}
}

// BeginTX begins a transaction on the test database with a rollback to be
// performed during the test cleanup
func (tdb *TestDB) BeginTX(t *testing.T) *sql.Tx {
	tx, err := tdb.DB.Begin()
	if err != nil {
		t.Fatalf("could not begin transaction: %v", err)
	}
	t.Cleanup(func() {
		tx.Rollback()
	})
	return tx
}

// NewTestRepo initializes a new test Repository with a transaction and
// populates it with sample rows
func NewTestRepo(t *testing.T) (repository.Repository, map[int64]types.RowDataMap) {
	tdb := NewTestDB(t)
	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx, tdb.Tables)
	sampleRows := InsertSampleRows(repo)
	return repo, sampleRows
}

// NewTestService initializes a new test Service with a test Repository
func NewTestService(t *testing.T) (service.Service, map[int64]types.RowDataMap) {
	repo, sampleRows := NewTestRepo(t)
	return service.NewService(repo), sampleRows
}
