package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"Auth-Service-Rest-Api/internal/handlers"
	"Auth-Service-Rest-Api/internal/db"
)

func main() {
	err := db.ConnectPostgresDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.ClosePostgresDB()
	db.WaitWhileDBNotReady()

	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/authorize", handlers.Authorize).Methods("POST")
	router.HandleFunc("/notes", handlers.CreateNote).Methods("POST")
	router.HandleFunc("/notes/{id}", handlers.GetNoteByID).Methods("GET")
	router.HandleFunc("/notes/{id}", handlers.UpdateNoteByID).Methods("PATCH")
	router.HandleFunc("/notes/{id}", handlers.DeleteNoteByID).Methods("DELETE")
	router.HandleFunc("/notes", handlers.ListNotes).Methods("GET")


	log.Fatal(http.ListenAndServe(":8080", router))
}
