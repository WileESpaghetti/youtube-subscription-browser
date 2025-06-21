CREATE TABLE IF NOT EXISTS video_tags (
    id INTEGER PRIMARY KEY,
    tag VARCHAR(255) NOT NULL,
    UNIQUE (tag) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS videos_video_tags (
    -- FIXME should there be a UNIQUE constraint on the IDs? Not an issue yet because we only import videos once
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     tag_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(tag_id) REFERENCES video_tags(id)
);
