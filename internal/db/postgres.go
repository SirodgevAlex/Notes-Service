package db

import (
	"database/sql"
	"log"
	"fmt"
	"time"
    "errors"
    "strconv"

    "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
    "Auth-Service-Rest-Api/internal/models"
    "Auth-Service-Rest-Api/internal/auth"
    "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var jwtKey = []byte("1234")

func ConnectPostgresDB() error {
    connStr := "postgres://postgres:1234@host.docker.internal:5432/notes_service?sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    err = db.Ping()
    if err != nil {
        return err
    }
    log.Println("Connected to PostgreSQL database")
    return nil
}

func ClosePostgresDB() {
    if db != nil {
        db.Close()
        log.Println("Disconnected from PostgreSQL database")
    }
}

func GetPostgresDB() (*sql.DB, error) {
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func WaitWhileDBNotReady() {
    fmt.Println("Waiting for database to be ready...")
    for {
        if err := db.Ping(); err == nil {
            fmt.Println("Database is ready!")
            break
        }
        fmt.Println("Database is not ready, waiting...")
        time.Sleep(time.Second)
    }
}

func RegisterUser(user models.User) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE Email = $1", user.Email).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
	}
	if count > 0 {
		return 0, errors.New("Email уже занят")
	}

	if !auth.IsEmailValid(user.Email) || !auth.IsPasswordSafe(user.Password) {
		return 0, errors.New("Некорректный email или пароль")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("ошибка при генерации пароля: %v", err)
	}

	query := "INSERT INTO users(Email, Password) VALUES($1, $2) RETURNING Id"
	err = db.QueryRow(query, user.Email, string(hashedPassword)).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при регистрации пользователя: %v", err)
	}

	return user.ID, nil
}

func AuthorizeUser(user models.User) (string, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT Password FROM Users WHERE Email = $1", user.Email).Scan(&hashedPassword)
	if err != nil {
		return "", fmt.Errorf("ошибка при поиске пользователя: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return "", errors.New("неправильный пароль")
	}

	var userID int
	err = db.QueryRow("SELECT Id FROM users WHERE Email = $1", user.Email).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("ошибка при получении ID пользователя: %v", err)
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &auth.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(userID),
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации токена: %v", err)
	}

	return tokenString, nil
}

func CreateNote(note *models.Note) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO notes(created_at, author_id, title, text) VALUES($1, $2, $3, $4) RETURNING id",
        note.CreatedAt.Format("2006-01-02 15:04:05"), note.AuthorID, note.Title, note.Text,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert note into database: %v", err)
	}
	if id == 0 {
		return 0, errors.New("failed to get ID of inserted note")
	}
	return id, nil
}

func UpdateNoteByID(id string, updatedNote *models.Note) error {
    query := `
        UPDATE notes
        SET 
            title = $1,
            text = $2
        WHERE id = $3
    `

    _, err := db.Exec(query, updatedNote.Title, updatedNote.Text, id)
    if err != nil {
        return fmt.Errorf("failed to update note: %v", err)
    }

    return nil
}

func GetNoteByID(id string) (*models.Note, error) {
	var note models.Note
	err := db.QueryRow("SELECT id, created_at, author_id, title, text FROM notes WHERE id = $1", id).Scan(
		&note.ID, &note.CreatedAt, &note.AuthorID, &note.Title, &note.Text,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("note not found")
	case err != nil:
		return nil, err
	}
	return &note, nil
}

func DeleteNoteByID(id string) error {
	result, err := db.Exec("DELETE FROM notes WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении заметки: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества затронутых строк: %v", err)
	}
	
	if rowsAffected == 0 {
		return errors.New("заметка с указанным ID не найдена")
	}

    return nil
}

func GetUserIDFromNote(id string) (string, error) {
    var authorID string
    err := db.QueryRow("SELECT author_id FROM notes WHERE id = $1", id).Scan(&authorID)
    if err != nil {
        return "", err
    }
    return authorID, nil
}


func GetCreatedAtFromNote(id string) (time.Time, error) {
    var createdAt time.Time
    err := db.QueryRow("SELECT created_at FROM notes WHERE id = $1", id).Scan(&createdAt)
    if err != nil {
        return time.Time{}, err
    }
    return createdAt, nil
}

func ListNotes(page int, startDate, endDate string) ([]models.Note, error) {
    offset := (page - 1) * 10

    query := "SELECT id, title, text, author_id FROM notes WHERE 1=1"
    if startDate != "" && endDate != "" {
        query += " AND created_at BETWEEN $1 AND $2"
    }
    query += " ORDER BY created_at DESC LIMIT 10 OFFSET $3"

    rows, err := db.Query(query, startDate, endDate, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notes []models.Note

    for rows.Next() {
        var note models.Note
        err := rows.Scan(&note.ID, &note.Title, &note.Text, &note.AuthorID)
        if err != nil {
            return nil, err
        }
        notes = append(notes, note)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return notes, nil
}