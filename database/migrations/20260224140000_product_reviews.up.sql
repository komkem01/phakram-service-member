SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS product_reviews (
  id uuid PRIMARY KEY,
  member_id uuid NOT NULL REFERENCES members (id),
  product_id uuid NOT NULL REFERENCES products (id),
  order_id uuid NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  order_item_id uuid NOT NULL UNIQUE REFERENCES order_items (id) ON DELETE CASCADE,
  rating smallint NOT NULL,
  comment text NOT NULL,
  is_visible boolean NOT NULL DEFAULT true,
  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  CONSTRAINT product_reviews_rating_check CHECK (rating >= 1 AND rating <= 5)
);

--bun:split

CREATE INDEX IF NOT EXISTS product_reviews_product_visible_created_idx
  ON product_reviews (product_id, is_visible, created_at DESC);

--bun:split

CREATE INDEX IF NOT EXISTS product_reviews_member_idx
  ON product_reviews (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS product_reviews_order_idx
  ON product_reviews (order_id);
