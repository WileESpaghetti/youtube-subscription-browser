CREATE TABLE IF NOT EXISTS keywords (
    id INTEGER PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    UNIQUE(keyword) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS channels_keywords (
    -- FIXME should there be a UNIQUE constraint on the IDs? Not an issue yet because we only imort channels once
    id INTEGER PRIMARY KEY,
    channel_id INTEGER,
    keyword_id INTEGER,
    FOREIGN KEY(channel_id) REFERENCES channels(id),
    FOREIGN KEY(keyword_id) REFERENCES keywords(id)
);
