package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/OmendraRathore/Content_Aggregation_System/internal/database"

	_ "github.com/lib/pq"
)

type appConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	queryStore := database.New(dbConn)

	app := appConfig{
		DB: queryStore,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", app.createUser)
	v1Router.Get("/users", app.withAuth(app.getUser))

	v1Router.Post("/feeds", app.withAuth(app.createFeed))
	v1Router.Get("/feeds", app.listFeeds)

	v1Router.Get("/feed_follows", app.withAuth(app.listFeedFollows))
	v1Router.Post("/feed_follows", app.withAuth(app.createFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", app.withAuth(app.deleteFeedFollow))

	v1Router.Get("/posts", app.withAuth(app.listPosts))

	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/err", errorHandler)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(queryStore, collectionConcurrency, collectionInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
