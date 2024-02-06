package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Mopcho/Golang-api/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not bound in environment")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("DB_URL not in environment variable")
	}

	conn, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	db := database.New(conn)
	apiCnfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

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

	v1Router.Get("/healthz", handle_readiness)
	v1Router.Get("/error", handleErr)
	v1Router.Post("/users", apiCnfg.handlerCreateUser)
	v1Router.Get("/users", apiCnfg.middlewareAuth(apiCnfg.handleGetUser))
	v1Router.Post("/feeds", apiCnfg.middlewareAuth(apiCnfg.handlerCreateFeed))
	v1Router.Post("/feeds_follows", apiCnfg.middlewareAuth(apiCnfg.handlerCreateFeedFollow))
	v1Router.Get("/feeds_follows", apiCnfg.middlewareAuth(apiCnfg.handlerGetFeedFollows))
	v1Router.Delete("/feeds_follows/{feedFollowID}", apiCnfg.middlewareAuth(apiCnfg.handlerDeleteFeedFollow))
	v1Router.Get("/feeds", apiCnfg.handleGetFeeds)

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
	fmt.Println("Port:", portString)
}
