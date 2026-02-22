DROP INDEX IF EXISTS idx_contact_messages_access_token;

ALTER TABLE contact_messages
  DROP COLUMN IF EXISTS access_token;
