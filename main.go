package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	db "github/gojogourav/RSSAggregator/db/sqlc"
	"log"
	"os"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT is not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database : ", err)
	}
	db := db.New(conn)

	apiConfig := apiConfig{
		DB: db,
	}

	go startScraping(
		db,
		10,
		time.Minute,
	)

	router := chi.NewRouter()

	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Post("/users", apiConfig.handlreCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))

	v1Router.Post("/post", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feed", apiConfig.handlerGetFeed)

	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollows))

	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
