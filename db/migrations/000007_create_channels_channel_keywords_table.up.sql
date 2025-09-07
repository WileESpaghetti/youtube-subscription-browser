CREATE TABLE IF NOT EXISTS channels_channel_keywords (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER,
    keyword_id INTEGER,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    FOREIGN KEY(keyword_id) REFERENCES keywords(id) ON DELETE CASCADE,
    UNIQUE(channel_id, keyword_id) ON CONFLICT IGNORE
);
