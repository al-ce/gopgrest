package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
)

var (
	dbname = os.Getenv("DB_NAME")
	dbpass = os.Getenv("DB_PASS")
)

func run() error {
	// Define connection params
	dbParams := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, dbpass, dbname,
	)

	// Open db connection
	db, err := sql.Open("postgres", dbParams)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Establish connection
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	APIHandler := NewAPIHandler(db)

	// Create server and routes
	mux := http.NewServeMux()
	mux.Handle("/", &APIHandler)

	// Run server
	log.Println("Listening on port 8090...")
	err = http.ListenAndServe(":8090", mux)
	log.Fatal(err)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type ExerciseSet struct {
	Name        string
	PerformedAt time.Time
	Weight      float32
	Unit        string
	Reps        int
	SetCount    int
	Notes       string
	SplitDay    string
	Program     string
	Tags        string
}
