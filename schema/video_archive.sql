-- WIP
-- FIXME some of these are from yt-dlp metadata which is not found in the api, maybe rename the extra stuff as archive_info?
CREATE TABLE IF NOT EXISTS video_archive_meta ( -- FIXME rename archive information?
                                                  id INTEGER PRIMARY KEY,
                                                  youtube_id TEXT NOT NULL,
                                                  title TEXT,
                                                  full_title TEXT,
                                                  description TEXT,
                                                  channel_id INTEGER,
                                                  width INTEGER,
                                                  height INTEGER,
                                                  resolution TEXT,
                                                  duration INTEGER,
                                                  webpage_url TEXT,
                                                  original_url TEXT,
                                                  uploaded_at INTEGER,
                                                  availability TEXT,
                                                  epoch INTEGER,
                                                  format TEXT,
                                                  format_id INTEGER,
                                                  format_note TEXT,
                                                  ext TEXT,
                                                  file_size INTEGER,
                                                  tbr REAL,
                                                  dynamic_range TEXT,
                                                  video_codec TEXT,
                                                  vbr REAL,
                                                  audio_codec TEXT,
                                                  aspect_ratio REAL,
                                                  abr REAL,
                                                  asr INTEGER,
                                                  status TEXT,

                                                  is_archived BOOLEAN DEFAULT FALSE,

                                                  FOREIGN KEY(channel_id) REFERENCES channels(id),
                                                  UNIQUE(youtube_id) ON CONFLICT IGNORE
);
