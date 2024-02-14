package main

import (
	"log"
	"net/http"
	"rss/internal/database"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apikeyString := strings.TrimPrefix(r.Header.Get("Authorization"), "ApiKey ")
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apikeyString)
		if err != nil {
			log.Print(err)
			respondWithError(w, 401, "Unauthorized!")
			return
		}
		handler(w, r, user)
	})
}
