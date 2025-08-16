package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"gopgrest/api"
	"gopgrest/repository"
	"gopgrest/service"
)

var (
	host       = os.Getenv("HOST")
	port       = os.Getenv("TEST_DB_PORT")
	user       = os.Getenv("TEST_DB_USER")
	password   = os.Getenv("TEST_DB_PASS")
	testDbName = os.Getenv("TEST_DB_NAME")

	AnneCarson = map[string]any{
		"born":     int64(1950),
		"died":     nil,
		"forename": "Anne",
		"id":       int64(1),
		"surname":  "Carson",
	}

	AnneBrontë = map[string]any{
		"born":     int64(1820),
		"died":     int64(1849),
		"forename": "Anne",
		"id":       int64(2),
		"surname":  "Brontë",
	}

	VirginiaWoolf = map[string]any{
		"born":     int64(1882),
		"died":     int64(1941),
		"forename": "Virginia",
		"id":       int64(3),
		"surname":  "Woolf",
	}
)

type TestDB struct {
	DB     *sql.DB
	TX     *sql.Tx
	Tables []repository.Table
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
func NewTestRepo(t *testing.T) repository.Repository {
	tdb := NewTestDB(t)
	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx, tdb.Tables)
	return repo
}

// NewTestService initializes a new test Service with a test Repository, using
// a transaction and returning the service plus some inserted sample rows
func NewTestService(t *testing.T) service.Service {
	repo := NewTestRepo(t)
	return service.NewService(repo)
}

// NewTestAPIHandler initializes an api handler with a transaction and return
// the handler plus some inserted sample rows
func NewTestAPIHandler(t *testing.T) api.APIHandler {
	tdb := NewTestDB(t)
	tx := tdb.BeginTX(t)
	h := api.NewAPIHandler(tx, tdb.Tables)
	return h
}
