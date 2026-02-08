SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_wishlist (
    id uuid PRIMARY KEY,
    member_id uuid REFERENCES members (id),
    product_id uuid REFERENCES products (id),
    quantity int,
    price_per_unit decimal,
    total_item_amount decimal,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_wishlist_member_id_idx ON member_wishlist (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_wishlist_product_id_idx ON member_wishlist (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_wishlist_member_product_uidx ON member_wishlist (member_id, product_id);
