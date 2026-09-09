package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pkab00/shortenit/internal/db"
	"github.com/pkab00/shortenit/internal/link"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := db.NewDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	link.RegisterRoutes(db, mux)

	http.ListenAndServe(":8080", mux)
}
