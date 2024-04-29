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

	user.ID, err = db.RegisterUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func GetNoteByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	noteID := params["id"]
	
	note, err := db.GetNoteByID(noteID)
	if err != nil {
		http.Error(w, "Failed to fetch note", http.StatusInternalServerError)
		return
	}
	
	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
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

    createdAt, err := db.GetCreatedAtFromNote(noteID)
    if err != nil {
        http.Error(w, "Failed to get note creation time", http.StatusInternalServerError)
        return
    }

    duration := time.Since(createdAt)
    days := int(duration.Hours() / 24)

    if days > 1 {
        http.Error(w, "Too late to update note", http.StatusForbidden)
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

    w.WriteHeader(http.StatusOK)
}

// func ListNotes(page int, startDate, endDate string) ([]models.Note, error) {
//     query := "SELECT id, title, text, author FROM notes WHERE 1=1"
//     if startDate != "" && endDate != "" {
//         query += " AND created_at BETWEEN $1 AND $2"
//     }
//     query += " ORDER BY created_at DESC LIMIT 10 OFFSET $3"

//     rows, err := db.Query(query, startDate, endDate, page*10)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()

//     var notes []models.Note

//     for rows.Next() {
//         var note models.Note
//         err := rows.Scan(&note.ID, &note.AuthorID, &note.Title, &note.Text)
//         if err != nil {
//             return nil, err
//         }
//         notes = append(notes, note)
//     }

//     if err := rows.Err(); err != nil {
//         return nil, err
//     }

//     return notes, nil
// }