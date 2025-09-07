CREATE TABLE IF NOT EXISTS videos_video_tags (
    id INTEGER PRIMARY KEY,
    video_id INTEGER,
    tag_id INTEGER,
    FOREIGN KEY(video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY(tag_id) REFERENCES video_tags(id) ON DELETE CASCADE,
    UNIQUE(video_id, tag_id) ON CONFLICT IGNORE
);
