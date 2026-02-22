CREATE TABLE IF NOT EXISTS contact_messages (
  id UUID PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(200) NOT NULL,
  subject VARCHAR(200) NOT NULL,
  message TEXT NOT NULL,
  send_status VARCHAR(20) NOT NULL DEFAULT 'pending',
  send_error TEXT,
  sent_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contact_messages_created_at ON contact_messages (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contact_messages_send_status ON contact_messages (send_status);
