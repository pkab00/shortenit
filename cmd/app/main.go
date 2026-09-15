package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/pkab00/shortenit/docs"
	"github.com/pkab00/shortenit/internal/db"
	"github.com/pkab00/shortenit/internal/link"
	"github.com/pkab00/shortenit/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title			ShortenIt API
// @version		1.0
// @description	REST API for URL shortening service
// @host			localhost:8080
// @basePath		/
func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@database:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)
	db, err := db.NewDB(url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	link.RegisterRoutes(db, mux)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", middleware.Logger(mux))
}
