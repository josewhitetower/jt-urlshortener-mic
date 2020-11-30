package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"urlshortener/models"

	// used to get the params from the route
	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"
)

// response format
type Response struct {
	OriginalURL string `json:"original_url"`
	ShortURL    int64  `json:"short_url"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file only for local env
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

// CreateURL create a user in the postgres db
func CreateURL(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var url models.URL

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// find existing url
	existingURL, err := getURLByOriginalURL(url.OriginalURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := Response{}
	if existingURL.OriginalURL != "" {
		// format a response object
		res.OriginalURL = existingURL.OriginalURL
		res.ShortURL = existingURL.ShortURL

	} else {
		// call insert user function and pass the user
		shortURL := insertURL(url)

		// format a response object
		res.OriginalURL = url.OriginalURL
		res.ShortURL = shortURL

		// send the response
	}
	json.NewEncoder(w).Encode(res)

}

// GetAllURLs will return all the urls
func GetAllURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the urls in the db
	urls, err := getAllURLs()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// send all the urls as response
	json.NewEncoder(w).Encode(urls)
}

// RedirectURL returns the orginal url to be redirected to
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the urls in the db
	shortURL := r.URL.Path[len("/api/shorturl/url/"):]
	n, err := strconv.ParseInt(shortURL, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	url, err := getURLByShortURL(n)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url.OriginalURL, 301)

}

//------------------------- handler functions ----------------
// insert one url in the DB
func insertURL(url models.URL) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning short_url will return the short_url of the inserted url
	sqlStatement := `INSERT INTO urls (original_url) VALUES ($1) RETURNING short_url`

	// the inserted shortURL will store in this shortURL
	var shortURL int64

	// execute the sql statement
	// Scan function will save the insert shortURL in the shortURL
	err := db.QueryRow(sqlStatement, url.OriginalURL).Scan(&shortURL)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", shortURL)

	// return the inserted id
	return shortURL
}

// get all urls
func getAllURLs() ([]models.URL, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var urls []models.URL

	// create the select sql query
	sqlStatement := `SELECT * FROM urls`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var url models.URL

		// unmarshal the row object to user
		err := rows.Scan(&url.OriginalURL, &url.ShortURL)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		urls = append(urls, url)

	}

	// return empty user on error
	return urls, err
}

func getURLByShortURL(shortURL int64) (models.URL, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a url of models.URL type
	var url models.URL

	// create the insert sql query
	// returning short_url will return the short_url of the inserted url
	sqlStatement := `SELECT * FROM urls WHERE short_url=$1`
	// execute the sql statement
	row := db.QueryRow(sqlStatement, shortURL)

	// unmarshal the row object to user
	err := row.Scan(&url.OriginalURL, &url.ShortURL)
	switch err {
	case sql.ErrNoRows:
		return url, errors.New("no rows were returned")
	case nil:
		return url, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	// return empty user on error
	return url, err
}

func getURLByOriginalURL(originalURL string) (models.URL, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a url of models.URL type
	var url models.URL

	// create the insert sql query
	sqlStatement := `SELECT * FROM urls WHERE original_url=$1`
	// execute the sql statement
	row := db.QueryRow(sqlStatement, originalURL)

	// unmarshal the row object to user
	err := row.Scan(&url.OriginalURL, &url.ShortURL)
	switch err {
	case sql.ErrNoRows:
		return url, nil
	case nil:
		return url, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	// return empty user on error
	return url, err
}
