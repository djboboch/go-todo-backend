package main

import (
	"log"
	"net/http"
)

func todoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Return all the to do's
	case "POST":
	// Create a new to do object and save it into DV
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func main() {

	http.HandleFunc("/api/v1/todo", todoHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
