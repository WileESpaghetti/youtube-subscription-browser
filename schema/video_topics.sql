CREATE TABLE IF NOT EXISTS video_topics (
    id INTEGER PRIMARY KEY,
    name TEXT,
    url TEXT,
    UNIQUE(url) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS videos_video_topics (
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     topic_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(topic_id) REFERENCES video_topics(id),
     UNIQUE(video_id, topic_id) ON CONFLICT IGNORE
);
