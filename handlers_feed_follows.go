package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rss/internal/database"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPostFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Something went wrong!")
		return
	}
	feedUUID, err := uuid.Parse(params.FeedID)
	if err != nil {
		log.Print(err)
		respondWithError(w, 401, "Unauthorized!")
		return
	}
	feed, err := cfg.DB.GetFeedByID(r.Context(), feedUUID)
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Something went wrong!")
		return
	}
	feed_follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	respondWithJSON(w, 200, feed_follow)
	return
}

// Maybe make it authenticated? Boot-dev says no
func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID := chi.URLParam(r, "id")
	feedFollowUUID, err := uuid.Parse(feedFollowID)
	if err != nil {
		log.Print(err)
		respondWithError(w, 401, "Unauthorized!")
		return
	}
	cfg.DB.DeleteFeedFollow(r.Context(), feedFollowUUID)
}

func (cfg *apiConfig) handlerGetFeedFollowByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := cfg.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		log.Print(err)
		respondWithError(w, 404, "Not found!")
		return
	}
	respondWithJSON(w, 200, feed_follows)
}
