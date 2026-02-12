SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS cart_items (
    id uuid PRIMARY KEY,
    cart_id uuid REFERENCES carts (id),
    product_id uuid REFERENCES products (id),
    quantity int,
    price_per_unit decimal,
    total_item_amount decimal,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS cart_items_cart_id_idx ON cart_items (cart_id);

--bun:split

CREATE INDEX IF NOT EXISTS cart_items_product_id_idx ON cart_items (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS cart_items_cart_product_uidx ON cart_items (cart_id, product_id);
