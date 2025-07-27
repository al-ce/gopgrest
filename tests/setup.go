package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
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
}

// ExerciseSet matches the exercise_set table in the test database so that rows
// can be scanned into fields of appropriate size
type ExerciseSet struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PerformedAt time.Time `json:"performed_at"`
	Weight      float32   `json:"weight"`
	Unit        string    `json:"unit"`
	Reps        int       `json:"reps"`
	SetCount    int       `json:"set_count"`
	Notes       string    `json:"notes"`
	SplitDay    string    `json:"split_day"`
	Program     string    `json:"program"`
	Tags        string    `json:"tags"`
}

// GetTestDB returns a test database
func GetTestDB(t *testing.T) *TestDB {
	testParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, testDbName,
	)

	db, err := sql.Open("postgres", testParams)
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return &TestDB{
		db,
		nil,
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
