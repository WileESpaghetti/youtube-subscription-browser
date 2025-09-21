CREATE TABLE IF NOT EXISTS archived_video_thumbnails (
    id INTEGER PRIMARY KEY,
    video_id INTEGER NOT NULL,
    resolution TEXT,
    index_id TEXT,
    preference INTEGER,
    width INTEGER,
    height INTEGER,
    url TEXT NOT NULL,
    file_name TEXT,
    FOREIGN KEY(video_id) REFERENCES videos(id) ON DELETE CASCADE
);