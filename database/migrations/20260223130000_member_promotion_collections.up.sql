CREATE TABLE IF NOT EXISTS member_promotion_collections (
  id UUID PRIMARY KEY,
  member_id UUID NOT NULL,
  promotion_id UUID NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
  collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(member_id, promotion_id)
);

CREATE INDEX IF NOT EXISTS idx_member_promotion_collections_member_id
  ON member_promotion_collections (member_id);
CREATE INDEX IF NOT EXISTS idx_member_promotion_collections_promotion_id
  ON member_promotion_collections (promotion_id);
