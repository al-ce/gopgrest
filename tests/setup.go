package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

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

func GetTestDB(t *testing.T) *TestDB {
	testParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, testDbName,
	)

	db, err := sql.Open("postgres", testParams)
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("could not begin transaction: %v", err)
	}

	t.Cleanup(func() {
		tx.Rollback()
		db.Close()
	})

	return &TestDB{
		db,
		tx,
	}
}
