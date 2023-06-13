package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kutelata/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// feed, err := urlToFeed("https://wagslane.dev/index.xml")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(feed)

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL_POSTGRES")
	if dbUrl == "" {
		log.Fatal("DB_URL_POSTGRES is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to the database", err)
	}

	apiConfig := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(apiConfig.DB, 10, time.Minute)

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
	v1Router.Get("/caolong", handleReadliness)
	v1Router.Get("/err", handleErr)

	v1Router.Post("/users", apiConfig.handleCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handleGetUser))

	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handleCreateFeed))
	v1Router.Get("/feeds", apiConfig.handleGetFeed)

	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handleGetPostsForUser))

	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handleDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
