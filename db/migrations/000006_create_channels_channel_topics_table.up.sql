CREATE TABLE IF NOT EXISTS channels_channel_topics (
    id INTEGER PRIMARY KEY,
    channel_id INTEGER,
    topic_id INTEGER,
    FOREIGN KEY(channel_id) REFERENCES channels(id) ON DELETE CASCADE,
    FOREIGN KEY(topic_id) REFERENCES topics(id) ON DELETE CASCADE,
    UNIQUE (channel_id, topic_id) ON CONFLICT IGNORE
);
