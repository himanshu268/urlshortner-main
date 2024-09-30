package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string `json:"id"`
	OriginalURL  string `json:"original_url"`
	ShortURL     string `json:"short_url"`
	CreationDate string `json:"creation_date"`
}

var urlDb = make(map[string]URL)

// Generates a short URL using MD5 hash
func generateshorturl(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

// Creates a new URL entry and stores it in the database
func createUrl(url string) string {
	shorturl := generateshorturl(url)

	urlDb[shorturl] = URL{
		ID:           shorturl,
		OriginalURL:  url,
		ShortURL:     shorturl,
		CreationDate: time.Now().Format(time.RFC3339), // Use a standard time format
	}

	return shorturl
}

// Retrieves the original URL from the database
func geturl(id string) (URL, error) {
	urlData, ok := urlDb[id]
	if !ok {
		return URL{}, errors.New("url not found")
	}
	return urlData, nil
}

// Simple hello world handler
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

// Handler for creating a short URL
func shorturlhandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OriginalURL string `json:"original_url"`
	}

	// Decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Create a short URL and respond with it
	shorturl := createUrl(data.OriginalURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shorturl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Redirect handler that redirects to the original URL
func redirect(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):] // Get the short URL ID4
	fmt.Println(id)
	url, err := geturl(id)
	fmt.Println((url)) // Retrieve original URL
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound) // Redirect to original URL
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/data", shorturlhandler) // Endpoint to create short URLs
	http.HandleFunc("/redirect/", redirect)   // Endpoint to redirect based on short URL
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
