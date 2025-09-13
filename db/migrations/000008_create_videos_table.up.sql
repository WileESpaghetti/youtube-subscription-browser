CREATE TABLE IF NOT EXISTS videos (
    id INTEGER PRIMARY KEY,
    youtube_id TEXT NOT NULL,
    published_at INTEGER,
    channel_id INTEGER,
    title TEXT
        COMMENT 'Trimmed to 100 characters by the API',
    description TEXT
        COMMENT 'Trimmed to 5000 bytes by the API',
    category_id TEXT,
    duration INTEGER,
    definition TEXT,
    is_licensed_content BOOLEAN,
    privacy_status TEXT,
    availability TEXT,
    media_type TEXT,
    webpage_url TEXT,
    original_url TEXT,
    view_count INTEGER,
    like_count INTEGER,
    dislike_count INTEGER,
    favorite_count INTEGER,
    comment_count INTEGER,
    is_archived BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    UNIQUE(youtube_id) ON CONFLICT REPLACE
);
