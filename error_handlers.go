package main

import "net/http"

// InternalServerErrorHandler responds with a 500 status and an error message
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

// NotFoundHandler responds with a 404 status and an error message
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
