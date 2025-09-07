CREATE TABLE IF NOT EXISTS videos_video_topics (
    id INTEGER PRIMARY KEY,
    video_id INTEGER,
    topic_id INTEGER,
    FOREIGN KEY(video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY(topic_id) REFERENCES video_topics(id) ON DELETE CASCADE,
    UNIQUE(video_id, topic_id) ON CONFLICT IGNORE
);
