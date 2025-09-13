CREATE TABLE IF NOT EXISTS channel_banners (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    UNIQUE(channel_id) ON CONFLICT IGNORE
)
