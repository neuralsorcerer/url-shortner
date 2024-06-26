package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }

    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer conn.Close(context.Background())

    db = conn
    fmt.Println("Connected to the database successfully")

    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler).Methods("GET")
    r.HandleFunc("/shorten", createShortURL).Methods("POST")
    r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")

    fmt.Println("Server is running at :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
