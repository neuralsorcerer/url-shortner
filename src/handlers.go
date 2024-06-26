package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"html/template"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func generateShortURL() string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, 6)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling createShortURL request")

    var url URL
    if err := r.ParseForm(); err != nil {
        log.Printf("Error parsing form: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    url.OriginalURL = r.FormValue("original_url")

    log.Printf("Original URL: %s", url.OriginalURL)

    shortURL := generateShortURL()
    url.ShortURL = shortURL

    log.Printf("Generated short URL: %s", url.ShortURL)

    _, err := db.Exec(context.Background(), "INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", url.OriginalURL, url.ShortURL)
    if err != nil {
        log.Printf("Error inserting into database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(url)
    log.Println("URL shortened successfully")
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortURL := vars["shortURL"]

    var originalURL string
    err := db.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE short_url=$1", shortURL).Scan(&originalURL)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, _ := template.ParseFiles("src/templates/index.html")
    tmpl.Execute(w, nil)
}
