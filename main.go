package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/WileESpaghetti/youtube-subscription-browser/api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

const port = ":8080"

func jsonError(w http.ResponseWriter, err interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func getAllChannels(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := api.ListResponse{}

		if r.Method != "GET" {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH405",
				Reason: "method not allowed",
			}
			jsonError(w, response, http.StatusMethodNotAllowed)
			return
		}

		channels, err := api.GetChannels(r.Context(), db)
		if err != nil {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH500",
				Reason: err.Error(),
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		total := len(channels)
		response.Page = api.Page{
			Page:         1,
			PerPage:      total,
			TotalPages:   1,
			TotalRecords: total,
		}

		// type conversion, probably want to use something fancier than []any in the future
		for _, c := range channels {
			response.Items = append(response.Items, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// initialize DB
	dbFile := "youtube.sqlite"
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFile))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// initialize HTTP server
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	http.Handle("/api/channels", getAllChannels(db))

	log.Printf("Listening on %s...", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
