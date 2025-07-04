CREATE TABLE IF NOT EXISTS video_formats (
    id INTEGER PRIMARY KEY,
    youtube_id TEXT NOT NULL, -- format_id: not enough info to tell if this should be unique or not
    video_id INTEGER,
    abr REAL,
    audio_codec TEXT,
    aspect_ratio REAL,
    asr INTEGER,
    audio_channels INTEGER,
    audio_ext TEXT,
    columns INTEGER,
    container TEXT,
    duration REAL,
    dynamic_range TEXT,
    ext TEXT,
    file_size INTEGER,
    file_size_approx INTEGER,
    format_note TEXT,
    fps REAL,
    has_drm BOOLEAN,
    height INTEGER,
    language TEXT,
    quality REAL,
    resolution TEXT,
    rows INTEGER,
    source_preference INTEGER,
    tbr REAL,
    url TEXT,
    vbr REAL,
    video_codec TEXT,
    video_ext TEXT,
    width INTEGER,
    requested BOOLEAN
    FOREIGN KEY(video_id) REFERENCES videos(id)
);
