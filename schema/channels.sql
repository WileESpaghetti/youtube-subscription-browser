CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY,
    youtube_id TEXT NOT NULL UNIQUE,
    title TEXT,
    description TEXT,
    custom_url TEXT,
    branding_title TEXT,
    branding_description TEXT,
    subscriber_count INTEGER NOT NULL default 0,
    video_count INTEGER NOT NULL default 0
);