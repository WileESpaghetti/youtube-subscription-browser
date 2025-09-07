CREATE TABLE IF NOT EXISTS videos_video_categories (
    id INTEGER PRIMARY KEY,
    video_id INTEGER,
    category_id INTEGER,
    FOREIGN KEY(video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY(category_id) REFERENCES video_categories(id) ON DELETE CASCADE,
    UNIQUE(video_id, category_id) ON CONFLICT IGNORE
);

