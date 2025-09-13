CREATE TABLE IF NOT EXISTS channel_thumbnails (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER NOT NULL,
    size TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    url TEXT NOT NULL,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    UNIQUE(channel_id, size) ON CONFLICT IGNORE
);
