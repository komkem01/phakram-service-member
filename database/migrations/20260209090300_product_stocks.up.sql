SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS product_stocks (
    id uuid PRIMARY KEY,
    product_id uuid REFERENCES products (id),
    unit_price decimal,
    stock_amount int,
    remaining int,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp,
    deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS product_stocks_product_id_idx ON product_stocks (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS product_stocks_product_uidx
    ON product_stocks (product_id)
    WHERE deleted_at IS NULL;
