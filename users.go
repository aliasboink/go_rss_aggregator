package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"rss/internal/database"
	"time"

	"github.com/google/uuid"
)

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (cfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Something went wrong!")
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Something went wrong!")
		return
	}
	respondWithJSON(w, 200, user)
	return
}
