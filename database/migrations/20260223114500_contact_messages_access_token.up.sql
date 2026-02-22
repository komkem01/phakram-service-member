ALTER TABLE contact_messages
  ADD COLUMN IF NOT EXISTS access_token VARCHAR(64);

UPDATE contact_messages
SET access_token = REPLACE(gen_random_uuid()::text, '-', '')
WHERE access_token IS NULL OR access_token = '';

ALTER TABLE contact_messages
  ALTER COLUMN access_token SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_messages_access_token
  ON contact_messages (access_token);
