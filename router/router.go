package router

import (
	"urlshortener/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/shorturl/url/{short_url}", middleware.RedirectURL).Methods("GET")
	router.HandleFunc("/api/shorturl/url", middleware.GetAllURLs).Methods("GET")
	router.HandleFunc("/api/shorturl/new", middleware.CreateURL).Methods("POST")

	return router
}
