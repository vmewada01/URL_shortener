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
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	CreatedDate  time.Time `json:"created_date"`
}

var urlDB = make(map[string]URL)

func generated_URL(original_url string) string {
	hasher := md5.New()
	hasher.Write([]byte(original_url))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func CreateDbUrl(original_url string) string {

	gen_url := generated_URL(original_url)
	id := gen_url
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  original_url,
		ShortenedURL: gen_url,
		CreatedDate:  time.Now(),
	}

	return gen_url
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func ShortenedURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	shortenUrl := CreateDbUrl(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortenUrl}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]

	url, err := getURL(id)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shortner", ShortenedURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)
	fmt.Println("starting at port 3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		errors.New("error in listening to the server")
	}
}
