package main

import (
	"log"
	"net/http"
	"rss/internal/database"
)

func (cfg *apiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostByUser(r.Context(), database.GetPostByUserParams{
		UserID: user.ID,
		Limit:  15,
	})
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Something went wrong!")
		return
	}
	respondWithJSON(w, 200, posts)
	return
}
