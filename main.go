package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"gopgrest/api"
	"gopgrest/repository"
)

func startServer() {
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

	// Establish connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tables, err := repository.GetPublicTables(db)
	if err != nil {
		panic(err)
	}
	APIHandler := api.NewAPIHandler(db, tables)

	// Create server and routes
	mux := http.NewServeMux()
	mux.Handle("/", &APIHandler)

	// Run server
	log.Printf("Listening on port %s...\n", apiport)
	http.ListenAndServe(":"+apiport, mux)
}

func main() {
	go startServer()

	// Block until quit signal
	quit := makeQuitListener()
	<-quit
	log.Println("Server is shutting down...")
}

// Make a channel to listen for a quit signal
func makeQuitListener() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func lookupEnv(varName string) string {
	varValue, exists := os.LookupEnv(varName)
	if !exists {
		panic(fmt.Sprintf("%s not found in environment", varName))
	}
	return varValue
}
