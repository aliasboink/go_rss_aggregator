package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"rss/internal/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment!")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL is not found in the environment!")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{
		DB: dbQueries,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1router.Post("/users", apiCfg.handlerPostUser)
	v1router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerPostFeed))
	v1router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerPostFeedFollow))
	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollowByUser))
	v1router.Delete("/feed_follows/{id}", apiCfg.handlerDeleteFeedFollow)
	v1router.Get("/readiness", handlerReadiness)
	v1router.Get("/err", handlerError)

	r.Mount("/v1", v1router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
