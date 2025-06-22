package api

import (
	"context"
	"database/sql"
	"strconv"
)

type Page struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type Error struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Reason string `json:"message"`
}

type ListResponse struct {
	// since this is single user and self-hosted we any pagination will be done from the frontend
	Items []any `json:"items"`
	Page  Page  `json:"meta"`
	Error Error `json:"error,omitempty"`
}

type Channel struct {
	ID                  int64  `json:"id"`
	YouTubeID           string `json:"youtube_id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	CustomURL           string `json:"custom_url"`
	BrandingTitle       string `json:"branding_title"`
	BrandingDescription string `json:"branding_description"`
	SubscriberCount     int64  `json:"subscriber_count"`
	VideoCount          int64  `json:"video_count"`
	IsArchived          bool   `json:"is_archived"`
}

func GetChannels(ctx context.Context, db *sql.DB) ([]Channel, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, youtube_id, title, description, custom_url, branding_title, branding_description, subscriber_count, video_count, is_archived FROM channels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		c := Channel{}

		if err := rows.Scan(
			&c.ID,
			&c.YouTubeID,
			&c.Title,
			&c.Description,
			&c.CustomURL,
			&c.BrandingTitle,
			&c.BrandingDescription,
			&c.SubscriberCount,
			&c.VideoCount,
			&c.IsArchived); err != nil {
			return nil, err
		}

		channels = append(channels, c)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func GetChannel(ctx context.Context, db *sql.DB, channelID string) (Channel, error) {
	var c Channel

	id, err := strconv.ParseInt(channelID, 10, 64)
	if err != nil {
		return c, err
	}

	err = db.QueryRowContext(ctx, "SELECT id, youtube_id, title, description, custom_url, branding_title, branding_description, subscriber_count, video_count, is_archived FROM channels WHERE id = ?", id).
		Scan(
			&c.ID,
			&c.YouTubeID,
			&c.Title,
			&c.Description,
			&c.CustomURL,
			&c.BrandingTitle,
			&c.BrandingDescription,
			&c.SubscriberCount,
			&c.VideoCount,
			&c.IsArchived)
	if err != nil {
		return c, err
	}

	return c, nil
}
