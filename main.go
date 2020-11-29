package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	fmt.Println("Connecting to DB")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/shorturl/new", getNewURL).Methods("POST")
	r.HandleFunc("/api/shorturl/{value}", redirectURL).Methods("GET")

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	port := getPort()

	log.Println("Server running in port: " + port)
	log.Fatal(http.ListenAndServe(port, r))

}

func getNewURL(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		OriginalURL int    `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	type Body struct {
		NewURL string `json:"newUrl"`
	}
	var b Body
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(b.NewURL)

	// check if url exist

	// if exists, give it back

	// if not, create a new entry and give it back

	res := Response{1, b.NewURL}
	json, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)

}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)["value"]
	if _, err := strconv.Atoi(v); err == nil {
		// find the entry with v as id
	}

	type Response struct {
		OriginalURL int    `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	http.Redirect(w, r, "https://google.com", 301)

}

// GetPort the Port from the environment so we can run on Heroku
func getPort() string {
	port := os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
