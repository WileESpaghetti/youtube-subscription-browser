CREATE TABLE IF NOT EXISTS videos_video_topics (
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     topic_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(topic_id) REFERENCES topics(id),
     UNIQUE(video_id, topic_id) ON CONFLICT IGNORE
);
