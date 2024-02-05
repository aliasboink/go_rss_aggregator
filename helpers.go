package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rss/internal/database"
	"time"

	"github.com/google/uuid"
)

// id UUID PRIMARY KEY,
// created_at TIMESTAMP NOT NULL,
// updated_at TIMESTAMP NOT NULL,
// title TEXT NOT NULL,
// url TEXT UNIQUE NOT NULL,
// description TEXT,
// published_at TIMESTAMP NOT NULL,
// feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	UserID      uuid.UUID
	PublishedAt string
	FeedID      uuid.UUID
}

func postToDatabasePost(post Post) (database.CreatePostParams, error) {
	desc := sql.NullString{Valid: false}
	if post.Description != "" {
		desc = sql.NullString{
			String: post.Description,
			Valid:  true,
		}
	}
	parsedTime, err := time.Parse(time.RFC1123Z, post.PublishedAt)
	if err != nil {
		return database.CreatePostParams{}, err
	}
	return database.CreatePostParams{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: desc,
		PublishedAt: parsedTime,
		FeedID:      post.FeedID,
	}, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
