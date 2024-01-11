package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/ScriptedPranav/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	// Importing postgres driver, mentioned in sqlc docs
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v\n",err)
	}

	portString := os.Getenv("PORT");
	if portString == "" {
		log.Fatal("Port not found")
	}

	dbURL := os.Getenv("DB_URL");
	if dbURL == "" {
		log.Fatal("Port not found")
	}

	conn,err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database",err)
	}
	
	//converts the sql.DB object to a database.Queries object
	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}
	
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins : []string{"https://*","http://*"},
		AllowedMethods : []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders : []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/user",apiCfg.handlerGetUser)

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server starting on %v",portString)

	er := srv.ListenAndServe()
	if er != nil {
		log.Fatal(er)
	}
}