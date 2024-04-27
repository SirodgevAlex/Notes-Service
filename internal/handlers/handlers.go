package handlers

import (
	"net/http"
	"encoding/json"
	"time"
	"strconv"

	"Auth-Service-Rest-Api/internal/models"
	"Auth-Service-Rest-Api/internal/auth"
	"Auth-Service-Rest-Api/internal/db"
	"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.RegisterUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := db.AuthorizeUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
    userID, err := auth.AuthenticateUser(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
    	return
    }

    var note models.Note
    if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
        http.Error(w, "Failed to decode JSON request body", http.StatusBadRequest)
        return
    }

    if len(note.Title) > 100 {
        http.Error(w, "Title is too long", http.StatusBadRequest)
        return
    }

    if len(note.Text) > 2000 {
        http.Error(w, "Text is too long", http.StatusBadRequest)
        return
    }

    note.CreatedAt = time.Now()
    note.AuthorID = userID

    noteID, err := db.CreateNote(&note)
    if err != nil {
        http.Error(w, "Failed to create note", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message": "Note created successfully",
        "note_id": noteID,
        "user_id": userID,
        "title":   note.Title,
        "text":    note.Text,
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
        return
    }
}

func UpdateNoteByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    noteID := params["id"]

    var updatedNote models.Note
    err := json.NewDecoder(r.Body).Decode(&updatedNote)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    note, err := db.GetNoteByID(noteID)
    if err != nil {
        http.Error(w, "Note not found", http.StatusNotFound)
        return
    }

    userID, err := auth.AuthenticateUser(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
    	return
    }

    noteAuthorID := note.AuthorID

    if userID != noteAuthorID {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    err = db.UpdateNoteByID(noteID, &updatedNote)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}


func DeleteNoteByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    noteID := params["id"]

    userID, err := db.GetUserIDFromNote(noteID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userIDFromToken, err := auth.AuthenticateUser(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
    	return
    }

    stringUserIDFromToken := strconv.Itoa(userIDFromToken)

    if userID != stringUserIDFromToken {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    err = db.DeleteNoteByID(noteID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
