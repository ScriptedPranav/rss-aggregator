package main

import (
	"log"
	"net/http"
	"os"
	"github.com/ScriptedPranav/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT");
	if portString == "" {
		log.Fatal("Port not found")
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

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server starting on %v",portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}