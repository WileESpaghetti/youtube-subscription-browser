package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

type ItemResponse struct {
	Item  any   `json:"item"`
	Error Error `json:"error,omitempty"`
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

type Video struct {
	ID          string `json:"id"`
	YouTubeID   string `json:"youtube_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ChannelID   string `json:"channel_id"`
	Duration    int64  `json:"duration"`
	WebpageURL  string `json:"webpage_url"`
	UploadedAt  int64  `json:"timestamp"`
	//Availability     string          `json:"availability"`
	OriginalURL string `json:"original_url"`
	FullTitle   string `json:"fulltitle"`
	//Epoch            int64           `json:"epoch"`
	//Format           string          `json:"format"`
	//FormatID         string          `json:"format_id"`
	//FormatNote       string          `json:"format_note"`
	//Ext              string          `json:"ext"`
	//FileSize         int64           `json:"filesize_approx"`
	//TBR              float32         `json:"tbr"`
	Width      int64  `json:"width"`
	Height     int64  `json:"height"`
	Resolution string `json:"resolution"`
	//DynamicRange     string          `json:"dynamic_range"`
	//VideoCodec       string          `json:"vcodec"`
	//VBR              float32         `json:"vbr"`
	//AudioCodec       string          `json:"acodec"`
	AspectRatio float32 `json:"aspect_ratio"` // < 1 = shorts/vert?
	//ABR              float32         `json:"abr"`
	//ASR              int64           `json:"asr"`
	//Categories       []string        `json:"categories"`
	//Tags             []string        `json:"tags"`
	//Formats          []YouTubeFormat `json:"formats"`
	//RequestedFormats []YouTubeFormat `json:"requested_formats"`
}

type ChannelVideoStats struct {
	ChannelID             int    `json:"channel_id"`
	ChannelTitle          string `json:"channel_title"`
	TotalVideos           int    `json:"total_videos"`
	TotalVideosArchived   int    `json:"total_videos_archived"`
	HasArchive            bool   `json:"has_archive"`
	HasCompleteArchive    bool   `json:"has_complete_archive"`
	LatestVideoUploadDate int    `json:"latest_video_upload_date"`
	LatestVideoYouTubeID  string `json:"latest_video_youtube_id"`
}

func GetChannelVideoStats(ctx context.Context, db *sql.DB, channelID int) (ChannelVideoStats, error) {
	var cvs ChannelVideoStats

	err := db.QueryRowContext(ctx, `
SELECT
    channels.id AS channel_id,
    channels.title,
    channels.video_count AS total_videos,
    COALESCE(archived_videos.archived_total, 0) AS total_videos_archived,
    (CASE WHEN archived_videos.archived_total > 0 THEN 1 ELSE 0 END) AS has_archive, -- FIXME need to check when no videos
    (CASE WHEN channels.video_count = archived_total THEN 1 ELSE 0 END) AS has_complete_archive,
    ranked_videos.uploaded_at AS lastest_video_upload_date,
    ranked_videos.youtube_id AS lastest_video_youtube_id
FROM channels
LEFT JOIN (SELECT * FROM (
         SELECT
             youtube_id,
             channel_id,
             uploaded_at,
             ROW_NUMBER() OVER (PARTITION BY channel_id ORDER BY uploaded_at DESC) as rn
         FROM videos
     )
WHERE rn = 1)
    AS ranked_videos ON channels.id = ranked_videos.channel_id
         LEFT JOIN (SELECT channel_id, COUNT(*) AS archived_total FROM videos GROUP BY videos.channel_id) archived_videos ON archived_videos.channel_id = channels.id
WHERE channels.id = ?
`, channelID).Scan(
		&cvs.ChannelID,
		&cvs.ChannelTitle,
		&cvs.TotalVideos,
		&cvs.TotalVideosArchived,
		&cvs.HasArchive,
		&cvs.HasCompleteArchive,
		&cvs.LatestVideoUploadDate,
		&cvs.LatestVideoYouTubeID)
	if err != nil {
		return cvs, err
	}

	return cvs, nil
}

func GetVideos(ctx context.Context, db *sql.DB, channelID int, fromTimestamp int) ([]Video, error) {
	whereClauses := make([]string, 0, 2)
	whereParams := make([]interface{}, 0, 2)

	if channelID > 0 {
		whereClauses = append(whereClauses, "channel_id = ?")
		whereParams = append(whereParams, channelID)
	}

	if fromTimestamp > 0 {
		whereClauses = append(whereClauses, "uploaded_at > ?")
		whereParams = append(whereParams, fromTimestamp)
	}

	stmt := "SELECT id, youtube_id, title, full_title, description, channel_id, width, height, resolution, duration, webpage_url, original_url, uploaded_at, aspect_ratio FROM videos" // mostly ignoring all of the format related fields
	if len(whereClauses) > 0 {
		stmt += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	fmt.Println(stmt)

	rows, err := db.QueryContext(ctx, stmt, whereParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		v := Video{}

		if err := rows.Scan(
			&v.ID,
			&v.YouTubeID,
			&v.Title,
			&v.FullTitle,
			&v.Description,
			&v.ChannelID,
			&v.Width,
			&v.Height,
			&v.Resolution,
			&v.Duration,
			&v.WebpageURL,
			&v.OriginalURL,
			&v.UploadedAt,
			&v.AspectRatio); err != nil {
			return nil, err
		}

		videos = append(videos, v)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
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
