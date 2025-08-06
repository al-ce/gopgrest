package test_utils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"gopgrest/api"
	"gopgrest/repository"
	"gopgrest/service"
	"gopgrest/types"
)

var (
	host       = os.Getenv("HOST")
	port       = os.Getenv("TEST_DB_PORT")
	user       = os.Getenv("TEST_DB_USER")
	password   = os.Getenv("TEST_DB_PASS")
	testDbName = os.Getenv("TEST_DB_NAME")
)

const TABLE1 = "exercise_sets"

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
func NewTestRepo(t *testing.T) (repository.Repository, types.RowDataIdMap) {
	tdb := NewTestDB(t)
	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx, tdb.Tables)
	sampleRows := InsertSampleRows(repo)
	return repo, sampleRows
}

// NewTestService initializes a new test Service with a test Repository, using
// a transaction and returning the service plus some inserted sample rows
func NewTestService(t *testing.T) (service.Service, types.RowDataIdMap) {
	repo, sampleRows := NewTestRepo(t)
	return service.NewService(repo), sampleRows
}

// NewTestAPIHandler initializes an api handler with a transaction and return
// the handler plus some inserted sample rows
func NewTestAPIHandler(t *testing.T) (api.APIHandler, types.RowDataIdMap) {
	tdb := NewTestDB(t)
	tx := tdb.BeginTX(t)
	apih := api.NewAPIHandler(tx, tdb.Tables)
	sampleRows := InsertSampleRows(apih.Repo)
	return apih, sampleRows
}
