package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/WileESpaghetti/youtube-subscription-browser/api"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

const port = ":8080"

func jsonError(w http.ResponseWriter, err interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err)
}

func getVideoStatsByChannelId(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := api.ListResponse{} // FIXME single item response

		if r.Method != "GET" {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH405",
				Reason: "method not allowed",
			}
			jsonError(w, response, http.StatusMethodNotAllowed)
			return
		}

		// channel_id parameter
		sChannelID := r.PathValue("id")
		channelID, err := strconv.Atoi(sChannelID)
		if len(sChannelID) != 0 && err != nil {
			response.Error = api.Error{
				Status: http.StatusBadRequest,
				Code:   "CH400",
				Reason: "channel_id is invalid",
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		c, err := api.GetChannelVideoStats(r.Context(), db, channelID)
		if err != nil {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH500",
				Reason: err.Error(),
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(c)

	}
}

func getAllVideos(db *sql.DB) http.HandlerFunc {
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

		// channel_id parameter
		sChannelID := r.URL.Query().Get("channel_id")
		channelID, err := strconv.Atoi(sChannelID)
		if len(sChannelID) != 0 && err != nil {
			response.Error = api.Error{
				Status: http.StatusBadRequest,
				Code:   "CH400",
				Reason: "channel_id is invalid",
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		sFrom := r.URL.Query().Get("from")
		from, err := strconv.Atoi(sFrom)
		if len(sFrom) != 0 && err != nil {
			response.Error = api.Error{
				Status: http.StatusBadRequest,
				Code:   "CH400",
				Reason: "from field is not a valid timestamp",
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		videos, err := api.GetVideos(r.Context(), db, channelID, from)
		if err != nil {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH500",
				Reason: err.Error(),
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		total := len(videos)
		response.Page = api.Page{
			Page:         1,
			PerPage:      total,
			TotalPages:   1,
			TotalRecords: total,
		}

		// type conversion, probably want to use something fancier than []any in the future
		for _, v := range videos {
			response.Items = append(response.Items, v)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}
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
		_ = json.NewEncoder(w).Encode(response)
	}
}

func getChannel(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := api.ListResponse{} // FIXME single item response

		if r.Method != "GET" {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH405",
				Reason: "method not allowed",
			}
			jsonError(w, response, http.StatusMethodNotAllowed)
			return
		}

		c, err := api.GetChannel(r.Context(), db, r.PathValue("id"))
		if err != nil {
			response.Error = api.Error{
				Status: http.StatusMethodNotAllowed,
				Code:   "CH500",
				Reason: err.Error(),
			}
			jsonError(w, response, http.StatusServiceUnavailable)
			return
		}

		resp := api.ItemResponse{Item: c}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
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
	http.Handle("/api/channels/{id}", getChannel(db))
	http.Handle("/api/channels/{id}/video_stats", getVideoStatsByChannelId(db))
	http.Handle("/api/videos", getAllVideos(db))

	log.Printf("Listening on %s...", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
