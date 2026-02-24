SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS product_review_images (
  id uuid PRIMARY KEY,
  review_id uuid NOT NULL REFERENCES product_reviews (id) ON DELETE CASCADE,
  image_url varchar NOT NULL,
  sort_order smallint NOT NULL DEFAULT 0,
  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS product_review_images_review_id_idx
  ON product_review_images (review_id);
