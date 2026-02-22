ALTER TABLE contact_messages
  DROP COLUMN IF EXISTS read_at,
  DROP COLUMN IF EXISTS is_read;
