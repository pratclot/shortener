// main.go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
)

var urlDatabase map[string]string

func init() {
	urlDatabase = make(map[string]string)
}

// ShortenHandler handles the shortening of URLs
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL, ok := requestBody["url"]
	if !ok || originalURL == "" {
		http.Error(w, "Missing URL in the request body", http.StatusBadRequest)
		return
	}

	shortID, err := shortid.Generate()
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	shortURL := "http://localhost:3000/" + shortID

	urlDatabase[shortID] = originalURL

	response := map[string]string{
		"originalUrl":   originalURL,
		"shortenedLink": shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RedirectHandler handles the redirection based on the short URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortID"]

	originalURL, ok := urlDatabase[shortID]
	if !ok {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/shorten", ShortenHandler).Methods("POST")
	r.HandleFunc("/{shortID}", RedirectHandler).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":3000", nil)
}
