CREATE TABLE IF NOT EXISTS video_thumbnails (
    id INTEGER PRIMARY KEY,
    video_id INTEGER NOT NULL,
    size TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    url TEXT NOT NULL,
    FOREIGN KEY(video_id) REFERENCES videos(id) ON DELETE CASCADE,
    UNIQUE (video_id, size) ON CONFLICT IGNORE
)
