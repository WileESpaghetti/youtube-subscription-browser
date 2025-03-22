-- https://developers.google.com/youtube/v3/docs/channels#topicDetails.topicIds[]
-- TODO wikipedia URL from topicDetails.topicCategories
CREATE TABLE IF NOT EXISTS topics(
    id INTEGER PRIMARY KEY,
    type VARCHAR(255), -- I'm OK with this not being normalized in a separate table
    topic_id VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS channels_topics (
id INTEGER PRIMARY KEY,
    channel_id INTEGER,
    topic_id INTEGER,
    UNIQUE (channel_id, topic_id) ON CONFLICT IGNORE,
    FOREIGN KEY(channel_id) REFERENCES channels(id),
    FOREIGN KEY(channel_id) REFERENCES channels(id)
);

-- Music
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('music', '/m/04rlf', 'Music (parent topic)'),
('music', '/m/02mscn', 'Christian music'),
('music', '/m/0ggq0m', 'Classical music'),
('music', '/m/01lyv', 'Country'),
('music', '/m/02lkt', 'Electronic music'),
('music', '/m/0glt670', 'Hip hop music'),
('music', '/m/05rwpb', 'Independent music'),
('music', '/m/03_d0', 'Jazz'),
('music', '/m/028sqc', 'Music of Asia'),
('music', '/m/0g293', 'Music of Latin America'),
('music', '/m/064t9', 'Pop music'),
('music', '/m/06cqb', 'Reggae'),
('music', '/m/06j6l', 'Rhythm and blues'),
('music', '/m/06by7', 'Rock music'),
('music', '/m/0gywn', 'Soul music');

-- Gaming
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('gaming', '/m/0bzvm2', 'Gaming (parent topic)'),
('gaming', '/m/025zzc', 'Action game'),
('gaming', '/m/02ntfj', 'Action-adventure game'),
('gaming', '/m/0b1vjn', 'Casual game'),
('gaming', '/m/02hygl', 'Music video game'),
('gaming', '/m/04q1x3q', 'Puzzle video game'),
('gaming', '/m/01sjng', 'Racing video game'),
('gaming', '/m/0403l3g', 'Role-playing video game'),
('gaming', '/m/021bp2', 'Simulation video game'),
('gaming', '/m/022dc6', 'Sports game'),
('gaming', '/m/03hf_rm', 'Strategy video game');

-- Sports
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('sports', '/m/06ntj 	', 'Sports (parent topic)'),
('sports', '/m/0jm_ 	', 'American football'),
('sports', '/m/018jz 	', 'Baseball'),
('sports', '/m/018w8 	', 'Basketball'),
('sports', '/m/01cgz 	', 'Boxing'),
('sports', '/m/09xp_ 	', 'Cricket'),
('sports', '/m/02vx4 	', 'Football'),
('sports', '/m/037hz 	', 'Golf'),
('sports', '/m/03tmr 	', 'Ice hockey'),
('sports', '/m/01h7lh 	', 'Mixed martial arts'),
('sports', '/m/0410tth 	', 'Motorsport'),
('sports', '/m/07bs0 	', 'Tennis'),
('sports', '/m/07_53 	', 'Volleyball');

-- Entertainment
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('entertainment ', '/m/02jjt', 'Entertainment (parent topic)'),
('entertainment ', '/m/09kqc', 'Humor'),
('entertainment ', '/m/02vxn', 'Movies'),
('entertainment ', '/m/05qjc', 'Performing arts'),
('entertainment ', '/m/066wd', 'Professional wrestling'),
('entertainment ', '/m/0f2f9', 'TV shows');

-- Lifestyle
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('lifestyle', '/m/019_rr', 'Lifestyle (parent topic)'),
('lifestyle', '/m/032tl', 'Fashion'),
('lifestyle', '/m/027x7n', 'Fitness'),
('lifestyle', '/m/02wbm', 'Food'),
('lifestyle', '/m/03glg', 'Hobby'),
('lifestyle', '/m/068hy', 'Pets'),
('lifestyle', '/m/041xxh', 'Physical attractiveness [Beauty]'),
('lifestyle', '/m/07c1v', 'Technology'),
('lifestyle', '/m/07bxq', 'Tourism'),
('lifestyle', '/m/07yv9', 'Vehicles');

-- Society
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('society', '/m/098wr', ' 	Society (parent topic)'),
('society', '/m/09s1f', ' 	Business'),
('society', '/m/0kt51', ' 	Health'),
('society', '/m/01h6rj', ' 	Military'),
('society', '/m/05qt0', ' 	Politics'),
('society', '/m/06bvp', ' 	Religion');

-- Other
INSERT OR IGNORE INTO topics(type, topic_id, description) VALUES
('other', '/m/01k8wb', 'Knowledge');
