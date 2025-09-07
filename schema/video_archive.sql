-- WIP
-- FIXME some of these are from yt-dlp metadata which is not found in the api, maybe rename the extra stuff as archive_info?
CREATE TABLE IF NOT EXISTS video_archive_meta ( -- FIXME rename archive information?
      id INTEGER PRIMARY KEY,
      youtube_id TEXT NOT NULL,
      title TEXT,
      full_title TEXT,
      description TEXT,
      channel_id INTEGER,

      webpage_url TEXT,
      original_url TEXT,
      uploaded_at INTEGER,
      availability TEXT,
      epoch INTEGER,
      status TEXT,









      format TEXT,
    -- format_id
      format_id INTEGER, -- format.youtube_id
    -- video_id
      abr REAL,
      audio_codec TEXT,
      aspect_ratio REAL,
      asr INTEGER,
    -- audio_channels
    -- audio_ext
    -- columns
    -- container
      duration INTEGER,
      dynamic_range TEXT,
      ext TEXT,
      file_size INTEGER,
    -- file_size_approx
      format_note TEXT,
    -- fps
    -- has_drm
      height INTEGER,
    -- language
    -- quality
      resolution TEXT,
    -- rows
    -- source_preference
      tbr REAL,
    -- url
      vbr REAL,
      video_codec TEXT,
    -- video_ext
      width INTEGER,
    -- requested

      FOREIGN KEY(channel_id) REFERENCES channels(id),
      UNIQUE(youtube_id) ON CONFLICT IGNORE
);
