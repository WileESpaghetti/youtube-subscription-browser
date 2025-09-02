CREATE TABLE IF NOT EXISTS video_tags (
    id INTEGER PRIMARY KEY,
    tag VARCHAR(255) NOT NULL,
    UNIQUE (tag) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS videos_video_tags (
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     tag_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(tag_id) REFERENCES video_tags(id),
     UNIQUE(video_id, tag_id) ON CONFLICT IGNORE
);
