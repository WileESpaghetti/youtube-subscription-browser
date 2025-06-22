package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type YouTubeFormat struct {
	// this mixes in image, audio, and video formats so we are not guaranteed to have all fields
	ABR              *float64 `json:"abr"`
	AudioCodec       *string  `json:"acodec"`
	AspectRatio      *float64 `json:"aspect_ratio"`
	ASR              *int64   `json:"asr"`
	AudioChannels    *int64   `json:"audio_channels"`
	AudioExt         *string  `json:"audio_ext"`
	Columns          *int64   `json:"columns"`
	Container        *string  `json:"container"`
	Duration         *float64 `json:"duration"`
	DynamicRange     *string  `json:"dynamic_range"`
	Ext              *string  `json:"ext"`
	FileSize         *int64   `json:"filesize"`
	FileSizeApprox   *int64   `json:"filesize_approx"`
	Format           *string  `json:"format"`
	YouTubeFormatID  *string  `json:"format_id"`
	FormatNote       *string  `json:"format_note"`
	FPS              *float64 `json:"fps"`
	HasDRM           *bool    `json:"has_drm"`
	Height           *int64   `json:"height"`
	Language         *string  `json:"language"`
	Quality          *float64 `json:"quality"`
	Resolution       *string  `json:"resolution"`
	Rows             *int64   `json:"rows"`
	SourcePreference *int64   `json:"source_preference"`
	TBR              *float64 `json:"tbr"`
	URL              *string  `json:"url"`
	VBR              *float64 `json:"vbr"`
	VideoCodec       *string  `json:"vcodec"`
	VideoExt         *string  `json:"video_ext"`
	Width            *int64   `json:"width"`
	Requested        *int64   `json:"requested"`
}

type Video struct {
	YouTubeID        string          `json:"id"`
	Title            string          `json:"title"`
	Type             string          `json:"_type"`
	Description      string          `json:"description"`
	ChannelID        string          `json:"channel_id"`
	Duration         int64           `json:"duration"`
	WebpageURL       string          `json:"webpage_url"`
	UploadedAt       int64           `json:"timestamp"`
	Availability     string          `json:"availability"`
	OriginalURL      string          `json:"original_url"`
	FullTitle        string          `json:"fulltitle"`
	Epoch            int64           `json:"epoch"`
	Format           string          `json:"format"`
	FormatID         string          `json:"format_id"`
	FormatNote       string          `json:"format_note"`
	Ext              string          `json:"ext"`
	FileSize         int64           `json:"filesize_approx"`
	TBR              float32         `json:"tbr"`
	Width            int64           `json:"width"`
	Height           int64           `json:"height"`
	Resolution       string          `json:"resolution"`
	DynamicRange     string          `json:"dynamic_range"`
	VideoCodec       string          `json:"vcodec"`
	VBR              float32         `json:"vbr"`
	AudioCodec       string          `json:"acodec"`
	AspectRatio      float32         `json:"aspect_ratio"` // < 1 = shorts/vert?
	ABR              float32         `json:"abr"`
	ASR              int64           `json:"asr"`
	Categories       []string        `json:"categories"`
	Tags             []string        `json:"tags"`
	Formats          []YouTubeFormat `json:"formats"`
	RequestedFormats []YouTubeFormat `json:"requested_formats"`
}

func saveFormats(ctx context.Context, db *sql.DB, videoID int, formats []YouTubeFormat, isRequested bool) error {
	for _, f := range formats {
		_, err := db.ExecContext(ctx, `INSERT INTO video_formats(
			youtube_id,
			video_id,
			abr,
			audio_codec,
			aspect_ratio,
			asr,
			audio_channels,
			audio_ext,
			columns,
			container,
			duration,
			dynamic_range,
			ext,
			file_size,
			file_size_approx,
			format_note,
			fps,
			has_drm,
			height,
			language,
			quality,
			resolution,
			rows,
			source_preference,
			tbr,
			url,
			vbr,
			video_codec,
			video_ext,
			width,
            requested)
			    VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			f.YouTubeFormatID, videoID, f.ABR, f.AudioCodec, f.AspectRatio, f.ASR, f.AudioChannels, f.AudioExt,
			f.Columns, f.Container, f.Duration, f.DynamicRange, f.Ext, f.FileSize, f.FileSizeApprox,
			f.FormatNote, f.FPS, f.HasDRM, f.Height, f.Language, f.Quality, f.Resolution,
			f.Rows, f.SourcePreference, f.TBR, f.URL, f.VBR, f.VideoCodec, f.VideoExt, f.Width, isRequested)
		if err != nil {
			fmt.Printf("...unable to save: %s\n", err)
			return err
		}
	}

	return nil
}

func saveTags(ctx context.Context, db *sql.DB, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	for _, t := range tags {
		result, err := db.ExecContext(ctx, "INSERT INTO video_tags(tag) VALUES(?)", t)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func getTagIDs(ctx context.Context, db *sql.DB, tags []string) ([]int, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	// handle IN clause placeholders
	tagPlaceholders := strings.Repeat("?,", len(tags))
	tagPlaceholders = tagPlaceholders[:len(tagPlaceholders)-1] // strip off the trailing ,
	args := make([]interface{}, 0, len(tags))
	for _, id := range tags {
		args = append(args, id)
	}

	queryTopicIDs := fmt.Sprintf("SELECT id FROM video_tags WHERE tag in (%s)", tagPlaceholders)
	rows, err := db.QueryContext(ctx, queryTopicIDs, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(tags))
	for rows.Next() {
		var tagID int
		if err := rows.Scan(&tagID); err != nil {
			return nil, err
		}

		ids = append(ids, tagID)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func saveTagAssociations(ctx context.Context, db *sql.DB, videoID int, tagIDs []int) error {
	// FIXME since keyword should be unique we could probably do this with an insert with select/join to automatically
	if len(tagIDs) == 0 {
		return nil
	}

	for _, tagID := range tagIDs {
		result, err := db.ExecContext(ctx, "INSERT INTO videos_video_tags(video_id, tag_id) VALUES(?, ?)", videoID, tagID)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func saveCategories(ctx context.Context, db *sql.DB, categories []string) error {
	if len(categories) == 0 {
		return nil
	}

	for _, c := range categories {
		result, err := db.ExecContext(ctx, "INSERT INTO video_categories(category) VALUES(?)", c)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func getCategoryIDs(ctx context.Context, db *sql.DB, categories []string) ([]int, error) {
	if len(categories) == 0 {
		return nil, nil
	}

	// handle IN clause placeholders
	categoryPlaceholders := strings.Repeat("?,", len(categories))
	categoryPlaceholders = categoryPlaceholders[:len(categoryPlaceholders)-1] // strip off the trailing ,
	args := make([]interface{}, 0, len(categories))
	for _, id := range categories {
		args = append(args, id)
	}

	queryTopicIDs := fmt.Sprintf("SELECT id FROM video_categories WHERE category in (%s)", categoryPlaceholders)
	rows, err := db.QueryContext(ctx, queryTopicIDs, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(categories))
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			return nil, err
		}

		ids = append(ids, categoryID)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func saveCategoryAssociations(ctx context.Context, db *sql.DB, videoID int, categoryIDs []int) error {
	// FIXME since keyword should be unique we could probably do this with an insert with select/join to automatically
	if len(categoryIDs) == 0 {
		return nil
	}

	for _, categoryID := range categoryIDs {
		result, err := db.ExecContext(ctx, "INSERT INTO videos_video_categories(video_id, category_id) VALUES(?, ?)", videoID, categoryID)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func getChannelID(db *sql.DB, youtubeID string) (int, error) {
	var id int

	err := db.QueryRow("SELECT id FROM channels WHERE youtube_id = ?", youtubeID).
		Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return id, err
	case err != nil:
		return id, err
	default:
		return id, nil
	}
}

func getVideoID(ctx context.Context, db *sql.DB, youtubeID string) (int, error) {
	var id int

	err := db.QueryRowContext(ctx, "SELECT id FROM videos WHERE youtube_id = ?", youtubeID).
		Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return id, err
	case err != nil:
		return id, err
	default:
		return id, nil
	}
}

func saveVideo(ctx context.Context, db *sql.DB, v *Video) error {
	fmt.Printf("%s: saving video...\n", v.YouTubeID)

	channelID, err := getChannelID(db, v.ChannelID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		err = fmt.Errorf(`%s: missing channel: channel ID = "%s"`, v.YouTubeID, v.ChannelID)
		fmt.Println(err)
		return err
	case err != nil:
		err = fmt.Errorf(`%s: missing channel: channel ID = "%s": unexpected errror %w`, v.YouTubeID, v.ChannelID, err)
		fmt.Println(err)
		return errors.Join()
	}

	_, err = db.ExecContext(ctx, `INSERT INTO videos(
		youtube_id,
		title,
		full_title,
		description,
		channel_id,
		width,
		height,
		resolution,
		duration,
		webpage_url,
		original_url,
		uploaded_at,
		availability,
		epoch,
		format,
		format_id,
		format_note,
		ext,
		file_size,
		tbr,
		dynamic_range,
		video_codec,
		vbr,
		audio_codec,
		aspect_ratio,
		abr,
		asr
) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		v.YouTubeID,
		v.Title,
		v.FullTitle,
		v.Description,
		channelID,
		v.Width,
		v.Height,
		v.Resolution,
		v.Duration,
		v.WebpageURL,
		v.OriginalURL,
		v.UploadedAt,
		v.Availability,
		v.Epoch,
		v.Format,
		v.FormatID,
		v.FormatNote,
		v.Ext,
		v.FileSize,
		v.TBR,
		v.DynamicRange,
		v.VideoCodec,
		v.VBR,
		v.AudioCodec,
		v.AspectRatio,
		v.ABR,
		v.ASR)
	if err != nil {
		err = fmt.Errorf(`%s: failed to save video: %w`, v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	videoID, err := getVideoID(ctx, db, v.YouTubeID)
	if err != nil {
		err = fmt.Errorf("%s: unable to get database ID: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	// Tags
	err = saveTags(ctx, db, v.Tags)
	if err != nil {
		err = fmt.Errorf("%s: unable to save tags: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	tagIDs, err := getTagIDs(ctx, db, v.Tags)
	if err != nil {
		err = fmt.Errorf("%s: unable to get tag IDs: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	err = saveTagAssociations(ctx, db, videoID, tagIDs)
	if err != nil {
		err = fmt.Errorf("%s: unable to assign tag to video: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	// Categories
	err = saveCategories(ctx, db, v.Categories)
	if err != nil {
		err = fmt.Errorf("%s: unable to save categories: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	categoryIDs, err := getCategoryIDs(ctx, db, v.Categories)
	if err != nil {
		err = fmt.Errorf("%s: unable to get category IDs: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	err = saveCategoryAssociations(ctx, db, videoID, categoryIDs)
	if err != nil {
		err = fmt.Errorf("%s: unable to assign categories to video: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	// Formats
	err = saveFormats(ctx, db, videoID, v.Formats, false)
	if err != nil {
		err = fmt.Errorf("%s: unable to save formats: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	err = saveFormats(ctx, db, videoID, v.RequestedFormats, true)
	if err != nil {
		err = fmt.Errorf("%s: unable to save requested formats: %w", v.YouTubeID, err)
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	if len(os.Args) <= 1 {
		// TODO show help
		_, _ = fmt.Fprintln(os.Stderr, "No input file specified")
		os.Exit(1)
	}

	ctx := context.Background()

	dbFile := "youtube.sqlite"
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFile))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, f := range os.Args[1:] {
		if strings.HasSuffix(f, ".json") {
			// TODO import single file
			continue
		}

		err := filepath.WalkDir(f, func(path string, d fs.DirEntry, err error) error {
			// TODO skip already imorted videos
			if d.IsDir() {
				return err
			}

			if filepath.Ext(d.Name()) != ".json" {
				return err
			}

			data, fileErr := os.Open(path)
			if err != nil {
				return errors.Join(err, fileErr)
			}
			defer data.Close()

			var o Video // consider data dirty because we could have random json files
			decoder := json.NewDecoder(data)
			decodeErr := decoder.Decode(&o)
			if err != nil {
				return errors.Join(err, decodeErr)
			}

			if o.Type != "video" {
				// I have playlist JSON files mixed in
				return err
			}

			videoErr := saveVideo(ctx, db, &o)
			if err != nil {
				return errors.Join(err, videoErr)
			}

			return err
		})

		fmt.Println("-----------------------------------------------")
		fmt.Printf("finished importing videos from: %s\n", f)
		fmt.Println(err)
	}
}
