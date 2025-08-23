CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY,
    youtube_id TEXT NOT NULL,
    title TEXT,
    description TEXT,
    custom_url TEXT,
    branding_title TEXT,
    branding_description TEXT,
    subscriber_count INTEGER NOT NULL default 0,
    video_count INTEGER NOT NULL default 0,
    uploads_playlist_id TEXT,
    UNIQUE(youtube_id) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS channel_thumbnails (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER NOT NULL,
    size TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    url TEXT NOT NULL,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE
)