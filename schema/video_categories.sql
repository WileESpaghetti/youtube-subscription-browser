CREATE TABLE IF NOT EXISTS video_categories (
    id INTEGER PRIMARY KEY,
    youtube_id VARCHAR(255) NOT NULL, -- These are numeric, but the API returns strings, so we'll store as a string to make lookups easier
    title VARCHAR(255) NOT NULL,
    assignable BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS videos_video_categories (
     id INTEGER PRIMARY KEY,
     video_id INTEGER,
     category_id INTEGER,
     FOREIGN KEY(video_id) REFERENCES videos(id),
     FOREIGN KEY(category_id) REFERENCES video_categories(id),
     UNIQUE(video_id, category_id) ON CONFLICT IGNORE
);

-- These come from the YouTube API
-- see: https://www.googleapis.com/youtube/v3/videoCategories/list?part=snippet&regionCode=US

INSERT INTO video_categories (youtube_id, title, assignable) VALUES
 ('1', 'Film & Animation', true),
('2', 'Autos & Vehicles', true),
('10', 'Music', true),
('15', 'Pets & Animals', true),
('17', 'Sports', true),
('18', 'Short Movies', false),
('19', 'Travel & Events', true),
('20', 'Gaming', true),
('21', 'Videoblogging', false),
('22', 'People & Blogs', true),
('23', 'Comedy', true),
('24', 'Entertainment', true),
('25', 'News & Politics', true),
('26', 'Howto & Style', true),
('27', 'Education', true),
('28', 'Science & Technology', true),
('29', 'Nonprofits & Activism', true),
('30', 'Movies', false),
('31', 'Anime/Animation', false),
('32', 'Action/Adventure', false),
('33', 'Classics', false),
('34', 'Comedy', false),
('35', 'Documentary', false),
('36', 'Drama', false),
('37', 'Family', false),
('38', 'Foreign', false),
('39', 'Horror', false),
('40', 'Sci-Fi/Fantasy', false),
('41', 'Thriller', false),
('42', 'Shorts', false),
('43', 'Shows', false),
('44', 'Trailers', false);