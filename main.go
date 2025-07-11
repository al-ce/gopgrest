package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
)

func run() error {
	// Define connection params
	dbname := lookupEnv("DB_NAME")
	dbpass := lookupEnv("DB_PASS")
	apiPort := lookupEnv("API_PORT")
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
	log.Printf("Listening on port %s...\n", apiPort)
	err = http.ListenAndServe(":"+apiPort, mux)
	log.Fatal(err)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func lookupEnv(varName string) string {
	varValue, exists := os.LookupEnv(varName)
	if !exists {
		panic(fmt.Sprintf("%s not found in environment", varName))
	}
	return varValue
}
