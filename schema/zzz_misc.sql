-- allow channels to be marked as archived
ALTER TABLE channels ADD COLUMN is_archived BOOLEAN NOT NULL DEFAULT 0;