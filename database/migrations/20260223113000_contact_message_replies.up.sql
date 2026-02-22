CREATE TABLE IF NOT EXISTS contact_message_replies (
  id UUID PRIMARY KEY,
  contact_message_id UUID NOT NULL REFERENCES contact_messages(id) ON DELETE CASCADE,
  sender_role VARCHAR(20) NOT NULL,
  sender_name VARCHAR(120) NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contact_message_replies_message_id_created_at
  ON contact_message_replies (contact_message_id, created_at ASC);
