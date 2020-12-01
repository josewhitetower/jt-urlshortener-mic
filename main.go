package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"urlshortener/router"
)

func main() {
	r := router.Router()

	port := getPort()

	log.Println("Server running in port: " + port)
	log.Fatal(http.ListenAndServe(port, r))

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
