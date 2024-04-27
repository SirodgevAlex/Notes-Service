package models

import "time"

type Note struct {
    ID           int       `json:"id"`
    CreatedAt    time.Time `json:"created_at"`
    AuthorID     int       `json:"author_id"`
	Title        string    `json:"title"`
    Text         string    `json:"text"`
}
