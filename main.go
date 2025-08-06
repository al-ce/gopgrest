package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"gopgrest/api"
	"gopgrest/repository"
)

func run() error {
	// Define connection params
	host := lookupEnv("HOST")
	dbuser := lookupEnv("DB_USER")
	dbname := lookupEnv("DB_NAME")
	dbpass := lookupEnv("DB_PASS")
	apiport := lookupEnv("API_PORT")
	dbport := lookupEnv("DB_PORT")
	dbparams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbport, dbuser, dbpass, dbname,
	)

	// Open db connection
	db, err := sql.Open("postgres", dbparams)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Establish connection
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	tables, err := repository.GetPublicTables(db)
	if err != nil {
		panic("Could not get public tables")
	}
	APIHandler := api.NewAPIHandler(db, tables)

	// Create server and routes
	mux := http.NewServeMux()
	mux.Handle("/", &APIHandler)

	// Run server
	log.Printf("Listening on port %s...\n", apiport)
	err = http.ListenAndServe(":"+apiport, mux)
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
