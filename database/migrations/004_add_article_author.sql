ALTER TABLE articles ADD COLUMN IF NOT EXISTS last_edited_by VARCHAR(255) DEFAULT 'system';

UPDATE articles SET last_edited_by = 'system' WHERE last_edited_by IS NULL;
