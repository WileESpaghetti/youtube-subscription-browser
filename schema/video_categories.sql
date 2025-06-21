CREATE TABLE IF NOT EXISTS video_categories (
    id INTEGER PRIMARY KEY,
    category VARCHAR(255) NOT NULL,
    UNIQUE (category) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS videos_video_categories (
    -- FIXME should there be a UNIQUE constraint on the IDs? Not an issue yet because we only import videos once
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     category_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(category_id) REFERENCES video_categories(id)
);
